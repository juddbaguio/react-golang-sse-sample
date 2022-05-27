# SSE-Golang and React Client Demo

## Event Source:
```zsh
http://localhost:3000/sse
```

### Notification Endpoint (POST):
```zsh
http://localhost:3000/notify
```
### Default message event (Sample Payload)
```json
{
    "data": {
        "alert": "You have a new notification!",
        "notificationMessage": "WOOOW"
    }
}
```

### Custom message event (Sample)
```json
{
    "event": "notify-me",
    "data": {
        "alert": "You have a new notification!",
        "notificationMessage": "WOOOW"
    }
}
```

## React Client

### Client runs at localhost:4000, you may change the PORT variable inside the .env file
------------------------
#### Change directory and install dependencies
```zsh
cd sse-client && npm i
or
cd sse-client && yarn
```

#### run the dev server
```zsh
npm run start
or 
yarn start
```
-------------------------
### SSE client logic
```javascript
function App() {
  const [sse, setSSE] = useState(null)
  useEffect(() => {
    // auto-reconnect logic of SSE
    if (sse) {
      // default message event
      sse.addEventListener("message", (evt) => {
          // rest of the code...
      })

      // custom message event
      sse.addEventListener("notify-me", evt => {
        const notificationObj = JSON.parse(evt.data)
        // rest of the code...
      })

      // onerror
      sse.onerror = () => {
        sse.close() // close connection always
        setSSE(new EventSource("http://localhost:3000/sse")) // reinstantiate SSE object
      }
    }
  }, [sse])

  // initial instantiation of opening Event Source
  useEffect(() => {
    setSSE(new EventSource("http://localhost:3000/sse"))
  }, [])

  return (
    // rest of the code...
  );
}
```