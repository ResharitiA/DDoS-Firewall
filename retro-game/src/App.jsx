import './App.css'

function App() {
  return (
    <div className="window">
      <div className="title-bar">
        <div className="title-bar-text">DDoS-Firewall</div>
        <div className="title-bar-controls">
          <button className="btn-close-red">X</button>
        </div>
      </div>

      <div className="window-body">
        {/* Главный слоган по брендбуку */}
        <h1 className="main-title">
          Надежная защита от DDoS-атак
        </h1>
        <p className="subtitle">
          DDoS-GUARD SECURITY INTERFACE v2.0
        </p>
        
        <div className="game-description">
          <p>
            <strong>ВНИМАНИЕ:</strong> Обнаружена сетевая угроза типа UDP-flood. <br/>
            Интеллектуальная система фильтрации <strong>DDoS-Guard</strong> готова к работе.
          </p>
          <hr className="divider" />
          <p className="footer-text">
            Геораспределенная сеть фильтрации активирована. 
            Все легитимные запросы проходят в штатном режиме.
          </p>
        </div>

        <button className="btn-play" onClick={() => alert('Протокол DDoS-Guard запущен!')}>
          ЗАПУСТИТЬ DDoS-FIREWALL
        </button>
      </div>
    </div>
  )
}

export default App