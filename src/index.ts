import 'dotenv/config'

import axios from 'axios'
import { verifyKeyMiddleware, InteractionType, InteractionResponseType } from 'discord-interactions'
import express, { Router } from 'express'

const app = express()

const routes = Router()

routes.get('/', (req, res) => res.status(200).json({ data: new Date().toUTCString() }))

routes.post('/interactions', verifyKeyMiddleware(String(process.env.PUBLIC_KEY)), async (req, res) => {
  const message = req.body

  if (message.type !== InteractionType.APPLICATION_COMMAND) {
    return res.status(400).json()
  }

  console.info('Message received', JSON.stringify(message))

  async function ping() {
    try {
      const url = message.data.options[0].value

      const response = await axios.get(url)

      return `\`\`\`json
${JSON.stringify(response.data)}
\`\`\`
      `
    } catch (err) {
      return `\`\`\`json
${JSON.stringify(err)}
\`\`\`
      `
    }
  }

  return res.status(200).json({
    type: InteractionResponseType.CHANNEL_MESSAGE_WITH_SOURCE,
    data: {
      embeds: [
        {
          color: 4437377,
          title: 'Pong',
          description: await ping(),
        },
      ],
    },
  })
})

app.use(String(process.env.BASE_PATH), routes)

app.listen(3000, () => console.log('It is running 🚀'))
