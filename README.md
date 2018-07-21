## Google Home[Base]

The idea behind Google Homebase is to allow Google Home to respond to the question:
> Is John at home?

**How does it do that?**

It uses iCloud to find an iPhone (currently hardcoded for an "iPhone X"), ideally it will use the phone ID. The service will compare the phone location (coordinates) to a predefined "homebase" location, and check if the phone is within about 25 meters from home.

This code can be deployed as aws lambda (to be exposed via Api Gateway as a POST endpoint) to provide this information to Google Home. The new `V2` interface is being used to marshal/unmarshal `dialogflow` request/response.

#### Configuration
This lambda has to be executed with a role that has access to KMS as the iCloud credentials are stored encrypted as env variables, therefore KMS access is needed in order to decrypt the credentials and issue the call to iCloud.

Environment variables needed for this to work:
- `UNAME`: iCloud username (encrypted with KMS)
- `PWD`: iCloud password (encrypted with KMS)
- `LAT`: float for the homebase latitude
- `LNG`: float for the homebase longitude