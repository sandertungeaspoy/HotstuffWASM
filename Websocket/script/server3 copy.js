var server = { urls: "stun:stun.l.google.com:19302" };
var dc3, pc3 = new RTCPeerConnection({ iceServers: [server] });
pc3.ondatachannel = e3 => dcInit3(dc3 = e3.channel);
pc3.oniceconnectionstatechange = e3 => console.log(pc3.iceConnectionState);


var offer3 = "";
var answer3 = "";
var welcomeMsg3 = "";

function dcInit3() {
    dc3.onopen = () => sendData3(welcomeMsg);
    dc3.onmessage = e => console.log(e.data);
}

function createOffer3() {
    dcInit3(dc3 = pc3.createDataChannel("chat"));
    pc3.createOffer().then(d => pc3.setLocalDescription(d));
    pc3.onicecandidate = e => {
        if (e.candidate) return;
        offer3 = pc3.localDescription.sdp;
        console.log(offer3);
    };
    return true;
};


function createAnswer3() {
    if (pc3.signalingState != "stable") return;
    var desc = new RTCSessionDescription({ type:"offer", sdp:offer3 });
    pc3.setRemoteDescription(desc)
        .then(() => pc3.createAnswer()).then(d => pc3.setLocalDescription(d));
};


function startChannel3() {
    console.log("Trying to start channel")
    if (pc3.signalingState != "have-local-offer") return;
    var desc = new RTCSessionDescription({ type:"answer", sdp:answer3 });
    pc3.setRemoteDescription(desc);
};

function sendData3(data) {
    dc3.send(data)
};

function startWebRTC3() {
    var id = document.getElementById("self-id").value;
    if (id == 1) {
        welcomeMsg3 = "Hello from Server 1"
        createOffer3();
        var ws = new WebSocket("ws://localhost:13374");
        ws.onopen = function() {
            offer3 += "&"
            ws.send(offer3);
            console.log(offer3);
            ws.send("setup:recvAnswer\n&");
            };
        
        ws.onmessage = function (evt) {
            answer3 = evt.data;
            console.log(answer3);
            startChannel3();
        }
    } else if (id == 4){
        // Wait for offer
        welcomeMsg3 = "Hello from Server 4"
        var ws = new WebSocket("ws://localhost:13374");
        ws.onopen = function() {
            ws.send("setup:recvOffer\n&");
            };
        ws.onmessage = function (evt) { 
            offer3 = evt.data;
            console.log(offer3);
            createAnswer3();
            
            pc3.onicecandidate = e => {
                if (e.candidate) return;
                answer3 = pc3.localDescription.sdp;
                console.log(answer3);
                answer3 += "&"; 
                ws.send(answer3);
            };
        };
    }
};
