# Cowin Telegram bot

Bot to lookup availability of vaccination slots in centers based on provided pincodes. if available then send message on telegram channel.
Follow this document for creating a telegram bot https://core.telegram.org/bots#6-botfather
Once you have bot token send a message to the bot from your telegram app and call below API to find the chat id which will be used by below app to send messages.
 `https://api.telegram.org/bot<YourBOTToken>/getUpdates`
  
More information available on this page https://stackoverflow.com/questions/32423837/telegram-bot-how-to-get-a-group-chat-id

Download appropriate release version from release page
https://github.com/codifierr/cowin_bot/releases

# Usage

Run the binary with below command arg. create a telegram bot using botfather and provide chat_id and telegram bot token and optional err_chat_id can be provided where bot will notify in case of any error while trying to get data from cowin apis
```
./cowin_bot -tel_token "" -chat_id "" -pincodes "" -min_age_limit "45" -min_available_capacity "4" -min_available_capacity_dose1 "4" -min_available_capacity_dose2 "0"
```
