import logo from './logo.svg';
import './App.css';
import { useEffect, useState } from 'react';
import { toast, ToastContainer } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';

function App() {
  const [sse, setSSE] = useState(null)
  useEffect(() => {
    if (sse) {
      // default message event
      sse.addEventListener("message", (evt) => {
        console.log(JSON.parse(evt.data))
      })

      // custom message event
      sse.addEventListener("notify-me", evt => {
        const notificationObj = JSON.parse(evt.data)
        toast(`🦄${notificationObj.notificationMessage}`, {
              position: "top-right",
              autoClose: 5000,
              hideProgressBar: false,
              closeOnClick: true,
              pauseOnHover: true,
              draggable: true,
              progress: undefined,
              })
      })

      // onerror
      sse.onerror = () => {
        sse.close()
        setSSE(new EventSource("http://localhost:3000/sse"))
      }
    }
  }, [sse])

  useEffect(() => {
    setSSE(new EventSource("http://localhost:3000/sse"))
  }, [])
  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        <p>
          Edit <code>src/App.js</code> and save to reload.
        </p>
        <a
          className="App-link"
          href="https://reactjs.org"
          target="_blank"
          rel="noopener noreferrer"
        >
          Learn React
        </a>
      </header>
      <ToastContainer />
    </div>
  );
}

export default App;
