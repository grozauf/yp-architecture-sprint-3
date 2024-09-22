#!/bin/sh

echo "### Test public management API..."

echo "#### Set deivce status -> "
curl --request POST \
  --url http://localhost:8000//manage/0/1/set \
  --header 'content-type: application/json' \
  --header 'user-agent: vscode-restclient' \
  --data '{"status": true,"value_name": "status"}'

echo "\n#### Set device target value -> "
curl --request POST \
  --url http://localhost:8000//manage/0/1/set \
  --header 'content-type: application/json' \
  --header 'user-agent: vscode-restclient' \
  --data '{"target_value": 1.01,"value_name": "target_value"}'

echo "\n#### Get device info -> "
curl --request GET \
  --url http://localhost:8000/manage/0/1/info \
  --header 'user-agent: vscode-restclient' \
  --data '{}'


echo "### Test public telemetry API..."

echo "\n#### Get latest telemetry value -> "
curl --request GET \
  --url http://localhost:8000/devices/0/1/telemetry/latest \
  --header 'user-agent: vscode-restclient' \
  --data '{}'

echo "\n#### Get history of telemetry values -> "
curl --request GET \
  --url http://localhost:8000/devices/0/1/telemetry \
  --header 'user-agent: vscode-restclient' \
  --data '{}'