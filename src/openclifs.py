import pandas as pd

import numpy as np
import json

prec = 5
df = pd.read_csv('spectrum.txt', skiprows=13, sep='\s+', decimal=',', header=None, names=['Wavelength', 'Counts']).round(prec)
dct = {'model_name': 'umuarama', 'wavelength': df.Wavelength.tolist()}

for i in range(1000):
    dct[f'sample_{i}'] = {
        "class": "soil", # this will change to soil, rock, etc
        "count": (df.Counts + np.random.rand()).round(prec).tolist()
    }

with open('saida.json', 'w+') as f:
    json.dump(dct, f)

# this is an output example:
# dct2 = {x: 'soil' for x in dct.keys() if 'count' in x}