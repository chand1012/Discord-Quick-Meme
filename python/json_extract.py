from json import loads

def json_extract(thing='', filename='data.json'):
    json_file = open(filename)
    json_data = loads(json_file.read())
    json_file.close()
    if not thing is '':
        return json_data[thing]
    else:
        return json_data
