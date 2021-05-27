# Cowin_bot

Bot to lookup availability of vaccination slots in centers based on provided pincodes. if available then send message on telegram channel.
Follow this document for creating a telegram bot https://docs.microsoft.com/en-us/azure/bot-service/bot-service-channel-connect-telegram?view=azure-bot-service-4.0

# Usage

Run the binary(mac os binary) which below command arg. create a telegram bot using botfather and provide chat_id and telegram bot token
```
./mac_locator -tel_token "" -chat_id "" -pincodes "" -min_age_limit "45" -min_available_capacity "4" -min_available_capacity_dose1 "4" -min_available_capacity_dose2 "0"
```
