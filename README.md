# GPT telegram bot
GPT Telegram Bot is a bot written in Golang that corrects your English grammar.

## Installation

- Create your own Telegram bot using BotFather and copy the API key into the .env file.
- Take your GPT API key from your [GPT account](https://platform.openai.com/examples) and put it into the .env file as well.
- Install [ngrock](https://ngrok.com/)
- Register for your webhook using this link https://api.telegram.org/bot{YOUR_TELEGRAM_API_TOKEN}/setWebhook?url=https://{YOUR_WEBHOOK_DOMAIN}/api/v1/update
- run docker-compose
```bash
docker-compose -f ./docker-compose.yaml up -d 
```
- Have fun!
