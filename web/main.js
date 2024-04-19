class CommandFromClient {
    constructor(heartbeatURL, RPS, targetRPSCheck) {
        this.HeartbeatUrl = heartbeatURL;
        this.RPS = RPS;
        this.TargetRPSCheck = targetRPSCheck;
    }
}

var ws;
var maxHeartbeats = 100;
var heartbeatData = [];
var heartbeatChart;

var maxLoadTestMarkers = 100;
var loadTestData = [];
var loadTestChart;

var startTime
var timeLastStateChanged
var stopwatchInterval
var elapsedPausedTime = 0; // to keep track of the elapsed time while stopped

var state = ""

window.addEventListener("load", function (evt) {
    document.getElementById("startform").onsubmit = function () {
        const command = new CommandFromClient(
            document.getElementById("heartbeatUrlInput").value,
            document.getElementById("loadRPS").value,
            document.getElementById("targetRPSTestCheck").checked);
        startTest(command);
        return false
    };

    heartbeatChartRender()
    loadTestChartRender()

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
            //addMessage("opened")
        }
        ws.onclose = function (evt) {
            //addMessage("closed")
            ws = null;
        }
        ws.onmessage = function (evt) {
            handleServerEvent(evt.data)
        }
        ws.onerror = function (evt) {
            handleServerEvent("error: " + evt.data)
        }
        return false;
    }
    else {
        console.log("Your browser does not support websockets")
    }

});

function handleServerEvent(msg) {
    try {
        const obj = JSON.parse(msg);
        if (obj.EventType == "heartbeat") {
            msg = epochMilisecondsToTime(obj.Data.Timestamp) + " - Ms Taken: " + obj.Data.MSLatency + " Msg: " + obj.Data.Message + " Success: " + obj.Data.Success
            updateHeartbeatChart(obj)
            addHeartbeatMessage(msg)
            updateState(obj.Data.Success)
        } else if (obj.EventType == "loadtest") {
            var vuString = obj.Data.VU.toString()
            if (vuString.length <= 1) {
                vuString = "0" + vuString
            }
            msg = vuString + " - " + epochMilisecondsToTime(obj.Data.Timestamp) + " - RPS: " + obj.Data.RPS
            updateLoadTestChart(obj)
            addLoadTestMessage(msg)
        }
    }
    catch (e) {
        console.log(e)
        addLoadTestMessage(msg)
    }
}

function updateState(success) {
    if (success) {
        if (state == "fail") {
            // We are in a state of fail, and we have changed to success, so we update blocktime
            var currentTime = new Date().getTime()
            var elapsedTimeInMiliseconds = currentTime - timeLastStateChanged
            console.log("time blocked for :" + elapsedTimeInMiliseconds)
            timeLastStateChanged = new Date().getTime()
        }
        state = "success"
    } else {
        if (state == "success") {
            // We have gone from success to fail
            var currentTime = new Date().getTime()
            var elapsedTimeInMiliseconds = currentTime - timeLastStateChanged
            console.log("time unblocked for :" + elapsedTimeInMiliseconds)
            timeLastStateChanged = new Date().getTime()
        }
        state = "fail"
    }
}

function updateHeartbeatChart(obj) {
    var dataPoint = []
    dataPoint.push(obj.Data.Timestamp)
    dataPoint.push(obj.Data.Success == true ? 1 : 0)

    if (heartbeatData.length >= maxHeartbeats) {
        heartbeatData.shift()
    }
    heartbeatData.push(dataPoint)

    heartbeatChart.updateSeries([{
        name: 'heartbeatChart',
        data: heartbeatData
    }])
}

function updateLoadTestChart(obj) {
    var dataPoint = []
    dataPoint.push(obj.Data.Timestamp)
    dataPoint.push(obj.Data.RPS)

    if (loadTestData.length >= maxLoadTestMarkers) {
        loadTestData.shift()
    }
    loadTestData.push(dataPoint)

    loadTestChart.updateSeries([{
        name: 'loadTestChart',
        data: loadTestData
    }])
}

function epochMilisecondsToTime(epoch) {
    var date = new Date(epoch)
    return date.toTimeString().split(" ")[0]
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
    startStopwatch()
    document.getElementById("startbutton").disabled = true;
    document.getElementById("stopbutton").disabled = false;
    document.getElementById("heartbeatUrlInput").disabled = true;
}

function stopTest() {
    if (!ws) {
        return
    }

    stopStopwatch()
    ws.send("stop");
    document.getElementById("startbutton").disabled = false;
    document.getElementById("stopbutton").disabled = true;
    document.getElementById("heartbeatUrlInput").disabled = false;
}

function scrollLogToBottom(t) {
    t.scrollTop = t.scrollHeight;
}

function heartbeatChartRender() {
    var options = {
        chart: {
            id: 'heartbeatChart',
            height: 250,
            type: 'line',
            animations: {
                enabled: false,
                easing: 'linear',
                dynamicAnimation: {
                    speed: 400
                }
            },
            toolbar: {
                show: false
            },
            zoom: {
                enabled: false
            }
        },
        dataLabels: {
            enabled: false
        },
        stroke: {
            curve: 'smooth'
        },
        markers: {
            size: 1
        },
        xaxis: {
            type: 'datetime'
        },
        series: [],
        title: {
            text: 'Heartbeat',
            align: 'left'
        },
        noData: {
            text: 'Loading...'
        },
        yaxis: {
            max: 1.5,
            min: 0
        },
        legend: {
            show: false
        },
    }

    heartbeatChart = new ApexCharts(
        document.querySelector("#heartbeatChart"),
        options
    );

    heartbeatChart.render();
}

function loadTestChartRender() {
    var options = {
        chart: {
            id: 'loadTestChart',
            height: 250,
            type: 'line',
            animations: {
                enabled: false,
                easing: 'linear',
                dynamicAnimation: {
                    speed: 400
                }
            },
            toolbar: {
                show: false
            },
            zoom: {
                enabled: false
            }
        },
        dataLabels: {
            enabled: false
        },
        stroke: {
            curve: 'smooth'
        },
        markers: {
            size: 1
        },
        xaxis: {
            type: 'datetime'
        },
        series: [],
        title: {
            text: 'rps',
            align: 'left'
        },
        noData: {
            text: 'Loading...'
        },
        legend: {
            show: false
        },
    }

    loadTestChart = new ApexCharts(
        document.querySelector("#loadTestChart"),
        options
    );

    loadTestChart.render();
}

function startStopwatch() {
    startTime = new Date().getTime()
    timeLastStateChanged = new Date().getTime()
    stopwatchInterval = setInterval(updateStopwatch, 50)
}

function updateStopwatch() {
    var currentTime = new Date().getTime()
    var elapsedTimeInMiliseconds = currentTime - startTime
    document.getElementById("stopwatch").innerHTML = milisecondsToFriendlyTime(elapsedTimeInMiliseconds)
    //console.log(elapsedTimeInMiliseconds)
    //var seconds = Math.floor(elapsedTimeInMiliseconds / 1000) % 60
    //console.log(seconds)
}

function milisecondsToFriendlyTime(miliseconds) {
    var milisecondsMod = Math.floor(miliseconds) % 1000
    var seconds = Math.floor(miliseconds / 1000) % 60
    var minutes = Math.floor(miliseconds / (1000 * 60)) % 60
    return String(minutes).padStart(2, '0') + ":" + String(seconds).padStart(2, '0') + "." + String(milisecondsMod).padStart(3, '0')
}

function stopStopwatch() {
    clearInterval(stopwatchInterval)
    stopwatchInterval = null
}

