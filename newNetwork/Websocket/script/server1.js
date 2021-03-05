// 192.168.50.73 for serving.
// 88.91.63.215 for websocket calls
var server = { urls: "stun:stun.l.google.com:19302" };

var dc, pc = new RTCPeerConnection({ iceServers: [server] });
pc.ondatachannel = e => dcInit(dc = e.channel);
// pc.oniceconnectionstatechange = e => console.log(pc.iceConnectionState);


var ws; 

var offer = "";
var answer = "";
var welcomeMsg = "";

function dcInit() {
    dc.onopen = () => sendData(welcomeMsg);
    dc.onmessage = e => console.log(e.data);
}

function createOffer() {
    dcInit(dc = pc.createDataChannel("chat"));
    pc.createOffer().then(d => pc.setLocalDescription(d));
    pc.addEventListener("icegatheringstatechange", ev => {
        switch(pc.iceGatheringState) {
          case "new":
            /* gathering is either just starting or has been reset */
            break;
          case "gathering":
            /* gathering has begun or is ongoing */
            break;
          case "complete":
            /* gathering has ended */
            offer = pc.localDescription.sdp;
            // console.log(offer);
            offer += "&"
            ws.send(offer);
            // console.log(offer);
            ws.send("setup:recvAnswer\n&");
            break;
        }
      });
    // pc.onicecandidate = e => {
        // console.log(e);
    //     if (e.candidate) return;
    //     offer = pc.localDescription.sdp;
        // console.log(offer);
    //     offer += "&"
    //     ws.send(offer);
        // console.log(offer);
    //     ws.send("setup:recvAnswer\n&");
    // };
    return true;
};


function createAnswer() {
    if (pc.signalingState != "stable") return;
    var desc = new RTCSessionDescription({ type:"offer", sdp:offer });
    pc.setRemoteDescription(desc)
        .then(() => pc.createAnswer()).then(d => pc.setLocalDescription(d));

    pc.addEventListener("icegatheringstatechange", ev => {
        switch(pc.iceGatheringState) {
            case "new":
            /* gathering is either just starting or has been reset */
            break;
            case "gathering":
            /* gathering has begun or is ongoing */
            break;
            case "complete":
            /* gathering has ended */
            if (pc.localDescription.sdp.includes("c=IN IP4 0.0.0.0")) {
                answer = pc.localDescription.sdp;
                // console.log(answer)
                answer += "&"; 
                ws.send(answer);
                ws.close();
            } else {
                // console.log("Retry");
                // console.log(pc.localDescription.sdp);
                restartWebRTC();
            };
        }
        });
    // pc.onicecandidate = e => {
        // console.log(e);
    //     if (e.candidate == "") {
            // console.log(answer);
    //         if (pc.localDescription.sdp.includes("c=IN IP4 0.0.0.0")) {
    //             answer = pc.localDescription.sdp;
    //             answer += "&"; 
    //             ws.send(answer);
    //             return;
    //         } else {
                // console.log("Retry");
    //             createAnswer();
    //         }
    //     }
    // };
};


function startChannel() {
    // console.log("Trying to start channel")
    if (pc.signalingState != "have-local-offer") return;
    var desc = new RTCSessionDescription({ type:"answer", sdp:answer });
    pc.setRemoteDescription(desc);
    ws.close();
};

function sendData(data) {
    dc.send(data)
};

function startWebRTC() {
    var id = document.getElementById("self-id").value;
    if (id == 1) {
        ws = new WebSocket("ws://localhost:13372");
        welcomeMsg = "Hello from Server 1"
        ws.onopen = function() {
            createOffer();
            // console.log("WS open");
        };
        ws.onmessage = function (evt) {
            answer = evt.data;
            // console.log(answer);
            startChannel();
            
        }
    } else if (id == 2) {
        // Wait for offer
        ws = new WebSocket("ws://localhost:13372");
        welcomeMsg = "Hello from Server 2"
        ws.onopen = function() {
            ws.send("setup:recvOffer\n&");
            };
        ws.onmessage = function (evt) { 
            offer = evt.data;
            // console.log(offer);
            createAnswer();   
        };
    }
};

function restartWebRTC() {
    dc, pc = new RTCPeerConnection({ iceServers: [server] });
    pc.ondatachannel = e => dcInit(dc = e.channel);
    // pc.oniceconnectionstatechange = e => console.log(pc.iceConnectionState);

    offer = "";
    answer = "";

    startWebRTC();
};
