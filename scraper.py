import pandas as pd
from pathlib import Path
from geopy.geocoders import Nominatim
from pprint import pprint
import re

geo_client = Nominatim(user_agent="Survivalrunbond-extended-calendar")

source = Path('./fixtures/wedstrijdschema.html')

with open(source, 'r') as f:
    dfs = pd.read_html(f.read(), header=0, extract_links='body', attrs={'class': 'wedstrijdagenda'})

# there is only one table of interest
df_wedstrijden = dfs[0]

# strip all tuple column names (pandas link extraction is a bit meh in 1.5.1)
# this breaks multispan columns though, so we're fixing those manually
#df_wedstrijden.columns = [t[0] for t in df_wedstrijden.columns ]
df_wedstrijden.columns = ['Datum', 'Plaats', 'LSR', 'MSR', 'KSR', 'Jeugd','BSR','Kwalificatie', 'Afstanden', 'min.leeftijd', 'Organisator', 'Inschrijflink', 'Uitslag']

# simply all column data that don't potentionally contain a link
for column in ['Datum', 'Plaats', 'LSR','MSR','KSR', 'Jeugd', 'BSR', 'Afstanden', 'min.leeftijd', 'Kwalificatie']:
    df_wedstrijden[column] = df_wedstrijden[column].apply(lambda t: t[0])

#retype klassement colums as boolean
for k in ['LSR', 'MSR', 'KSR', 'BSR', 'Jeugd', 'Kwalificatie']:
    df_wedstrijden[k] = df_wedstrijden[k].apply(lambda v: len(v) > 0)

#TODO: remove rows for nonsense events (e.g. 'ALV')

#parse available distances -> array
df_wedstrijden['Afstanden'] = df_wedstrijden['Afstanden'].apply(lambda a:  [n.replace(',','.') for n in re.findall('[\d,]+', a ) ])

#extract minimum age as proper number
df_wedstrijden['min.leeftijd'] = df_wedstrijden['min.leeftijd'].apply(lambda l:  next(iter(re.findall('\d+', l )), None) )

#ONK (open-nederlands-kampioenschap) -> new boolean
df_wedstrijden['ONK'] = df_wedstrijden['Plaats'].apply(lambda p: 'ONK' in p)

#Koppel -> new boolean
df_wedstrijden['Koppel'] = df_wedstrijden['Plaats'].apply(lambda p: 'Koppel' in p)

#afgelast -> boolean
df_wedstrijden['Afgelast'] = df_wedstrijden['Plaats'].apply(lambda p: 'AFGELAST' in p)

#cleanup placenames (remove content between brackets to indicate day-of-the-week, etc)
def cleanup_placename(p):
    mapping = {'(za)': '',
               '(zo)': '',
               '(B)': ', België',
               'ALV': '',
               'ONK': '',
               'Koppel': '',
               'BK': '',
               'KSR': '',
               'MSR': '',
               'JSR': '',
               'NSK': '',
               'LSR': '',
               'AFGELAST': '',
               '(gld)': 'Gelderland'
    }
    for (old, new) in mapping.items():
        p = p.replace(old, new)

    return p.strip()
df_wedstrijden['Plaats'] = df_wedstrijden['Plaats'].apply(cleanup_placename)

#lookup geocoordinates for each town
def lookup_coords(plaats):
    enriched_plaats = plaats
    if not "België" in plaats:
        enriched_plaats = plaats + ', Nederland'

    location = geo_client.geocode(enriched_plaats).raw
    return (location['lat'], location['lon'])

df_wedstrijden['coords'] = df_wedstrijden['Plaats'].apply(lookup_coords)

#print(df_wedstrijden.columns)
#print(df_wedstrijden[['Datum', 'Plaats', 'BSR']])

df_wedstrijden.to_json("wedstrijden.json", orient='records', indent=2)