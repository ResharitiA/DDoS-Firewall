import { useState, useEffect } from 'react'
import './App.css'

function App() {
  const [isGameRunning, setIsGameRunning] = useState(false);

  // Слушаем сообщение от iframe, чтобы закрыть его
  useEffect(() => {
    const handleMessage = (event) => {
      if (event.data === 'closeGame') {
        setIsGameRunning(false);
      }
    };

    window.addEventListener('message', handleMessage);
    return () => window.removeEventListener('message', handleMessage);
  }, []);

  if (isGameRunning) {
    return (
      <div className="game-wrapper">
        <button className="exit-button" onClick={() => setIsGameRunning(false)}>
          ✕ ВЕРНУТЬСЯ В ПАНЕЛЬ
        </button>
        <iframe 
          src="/gg/Cyber_defense.html" 
          title="Cyber Defense"
          className="game-iframe"
        />
      </div>
    );
  }

  return (
    <div className="window">
      {/* ... ваш остальной код окна (title-bar, window-body) ... */}
      <div className="title-bar">
        <div className="title-bar-text">DDoS-Firewall</div>
        <div className="title-bar-controls">
          <button className="btn-close-red">X</button>
        </div>
      </div>

      <div className="window-body">
        <h1 className="main-title">Надежная защита от DDoS-атак</h1>
        <p className="subtitle">DDoS-GUARD SECURITY INTERFACE v2.0</p>
        
        <div className="game-description">
          <p>
            <strong>ВНИМАНИЕ:</strong> Обнаружена сетевая угроза типа UDP-flood. <br/>
            Интеллектуальная система фильтрации <strong>DDoS-Guard</strong> готова к работе.
          </p>
          <hr className="divider" />
        </div>

        <button className="btn-play" onClick={() => setIsGameRunning(true)}>
          ЗАПУСТИТЬ DDoS-FIREWALL
        </button>
      </div>
    </div>
  )
}

export default App