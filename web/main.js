class CommandFromClient {
    constructor(heartbeatURL, RPS) {
        this.HeartbeatUrl = heartbeatURL;
        this.RPS = RPS;
    }
}

var ws;

window.addEventListener("load", function (evt) {
    document.getElementById("startform").onsubmit = function () {
        const command = new CommandFromClient(document.getElementById("heartbeatUrlInput").value,
            document.getElementById("loadRPS").value);
        startTest(command);
        return false
    };

    document.getElementById("stopbutton").onclick = function () {
        stopTest();
        return false
    };

    if (window["WebSocket"]) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("ws://" + window.location.hostname + ":8080/ws");
        ws.onopen = function (evt) {
            addMessage("opened")
        }
        ws.onclose = function (evt) {
            addMessage("closed")
            ws = null;
        }
        ws.onmessage = function (evt) {
            addMessage(evt.data)
        }
        ws.onerror = function (evt) {
            addMessage("error: " + evt.data)
        }
        return false;
    }
    else {
        addMessage("Your browser does not support websockets")
    }

});

function addMessage(msg) {
    var mylogarea = document.getElementById('mylog');
    var newValue = mylogarea.value + "\n" + msg;
    mylogarea.value = newValue
    scrollLogToBottom(mylogarea)
}

function startTest(commandFromClient) {
    if (!ws) {
        return
    }

    const myJSON = JSON.stringify(commandFromClient);

    ws.send("start " + myJSON);
}

function stopTest() {
    if (!ws) {
        return
    }

    ws.send("stop");
}

function scrollLogToBottom(t) {
    t.scrollTop = t.scrollHeight;
}

