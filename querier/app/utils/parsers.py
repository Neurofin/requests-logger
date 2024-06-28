import json

def textjson_to_json(json_text):
  json_start_index = json_text.find('{')
  json_end_index = json_text.rfind('}')

  json_text = json_text[json_start_index:json_end_index+1]
  jsonData = json.loads(json_text)
  return jsonData
