document.addEventListener('DOMContentLoaded', (event) => {
  let socket = new WebSocket('ws://localhost:8080/ws');
  socket.onmessage = (event) => {
    proccessWsData(event.data);
  };

  const form = document.getElementById('myForm');
  form.addEventListener('submit', sendData);
  let ul = document.getElementById('comments');
  let sents = document.getElementById('sents');

  function sendData(evt) {
    evt.preventDefault();
    const data = new FormData(evt.target);
    const dataObject = Object.fromEntries(data.entries());
    socket.send(dataObject['text']);
    form.reset();
  }

  const processChatMessage = (msg) => {
    let el = document.createElement('li');
    el.classList.add('has-text-right', 'has-text-primary');
    el.innerHTML = msg.data.Msg;
    ul.prepend(el);
  };

  const processSentiment = (msg) => {
    let el = document.createElement('li');
    el.classList.add(
      'has-text-right',
      'is-size-1',
      'animate__animated',
      'animate__fadeIn',
    );
    el.innerHTML = msg.data;
    sents.prepend(el);

    setTimeout(() => {
      sents.removeChild(el);
    }, 3000);
  };

  function proccessWsData(data) {
    const fnMap = {
      'ws.text.created': processChatMessage,
      'master.text.analyzed': processSentiment,
    };
    const msg = JSON.parse(data);
    const fn =
      fnMap[msg.type] ||
      function (msg) {
        console.log(msg);
      };

    fn(msg);
  }

  runSimulator();
});
