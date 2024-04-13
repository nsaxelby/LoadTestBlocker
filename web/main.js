class CommandFromClient {
    constructor(heartbeatURL, RPS, targetRPSCheck) {
        this.HeartbeatUrl = heartbeatURL;
        this.RPS = RPS;
        this.TargetRPSCheck = targetRPSCheck;
    }
}

var ws;

window.addEventListener("load", function (evt) {
    document.getElementById("startform").onsubmit = function () {
        const command = new CommandFromClient(
            document.getElementById("heartbeatUrlInput").value,
            document.getElementById("loadRPS").value,
            document.getElementById("targetRPSTestCheck").checked);
        startTest(command);
        return false
    };

    document.getElementById("stopbutton").onclick = function () {
        stopTest();
        return false
    };

    const checkbox = document.getElementById('targetRPSTestCheck')

    checkbox.addEventListener('change', (event) => {
        if (event.currentTarget.checked) {
            this.document.getElementById('loadRPS').disabled = false
        } else {
            this.document.getElementById('loadRPS').disabled = true
        }
    })

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
    if (msg.startsWith("Heartbeat:")) {
        msg = msg.replace("Heartbeat:", "")
        addHeartbeatMessage(msg)
    } else {
        addLoadTestMessage(msg)
    }
}

function addHeartbeatMessage(msg) {
    var heartbeatlogarea = document.getElementById('heartbeatoutputlog');
    var newValue = heartbeatlogarea.value + "\n" + msg;
    heartbeatlogarea.value = newValue
    scrollLogToBottom(heartbeatlogarea)
}

function addLoadTestMessage(msg) {
    var loadtestlogarea = document.getElementById('loadtestoutput');
    var newValue = loadtestlogarea.value + "\n" + msg;
    loadtestlogarea.value = newValue
    scrollLogToBottom(loadtestlogarea)
}

function startTest(commandFromClient) {
    if (!ws) {
        return
    }

    const myJSON = JSON.stringify(commandFromClient);

    ws.send("start " + myJSON);
    document.getElementById("startbutton").disabled = true;
    document.getElementById("stopbutton").disabled = false;
    document.getElementById("heartbeatUrlInput").disabled = true;
}

function stopTest() {
    if (!ws) {
        return
    }

    ws.send("stop");
    document.getElementById("startbutton").disabled = false;
    document.getElementById("stopbutton").disabled = true;
    document.getElementById("heartbeatUrlInput").disabled = false;
}

function scrollLogToBottom(t) {
    t.scrollTop = t.scrollHeight;
}
