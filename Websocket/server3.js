var server = { urls: "stun:stun.l.google.com:19302" };

var dc3, pc3 = new RTCPeerConnection({ iceServers: [server] });
pc3.ondatachannel = e3 => dcInit3(dc3 = e3.channel);
pc3.oniceconnectionstatechange = e3 => console.log(pc3.iceConnectionState);


var ws3;

var offer3 = "";
var answer3 = "";
var welcomeMsg3 = "";

function dcInit3() {
    dc3.onopen = () => sendData3(welcomeMsg3);
    dc3.onmessage = e => console.log(e.data);
}

function createOffer3() {
    dcInit3(dc3 = pc3.createDataChannel("chat"));
    pc3.createOffer().then(d => pc3.setLocalDescription(d));
    pc3.onicecandidate = e3 => {
    pc3.addEventListener("icegatheringstatechange", ev => {
        switch(pc3.iceGatheringState) {
            case "new":
            /* gathering is either just starting or has been reset */
            break;
            case "gathering":
            /* gathering has begun or is ongoing */
            break;
            case "complete":
            /* gathering has ended */
            offer3 = pc3.localDescription.sdp;
            console.log(offer3);
            offer3 += "&"
            ws3.send(offer3);
            console.log(offer3);
            ws3.send("setup:recvAnswer\n&");
            break;
        }
        });
        // console.log(e3);
        // if (e3.candidate) return;
        // offer3 = pc3.localDescription.sdp;
        // console.log(offer3);
        // offer3 += "&";
        // ws3.send(offer3);
        // console.log(offer3);
        // ws3.send("setup:recvAnswer\n&");
    };
    return true;
};


function createAnswer3() {
    if (pc3.signalingState != "stable") return;
    var desc = new RTCSessionDescription({ type:"offer", sdp:offer3 });
    pc3.setRemoteDescription(desc)
        .then(() => pc3.createAnswer()).then(d => pc3.setLocalDescription(d));
    
    pc3.addEventListener("icegatheringstatechange", ev => {
        switch(pc3.iceGatheringState) {
            case "new":
            /* gathering is either just starting or has been reset */
            break;
            case "gathering":
            /* gathering has begun or is ongoing */
            break;
            case "complete":
            /* gathering has ended */
            if (pc3.localDescription.sdp.includes("c=IN IP4 0.0.0.0")) {
                answer3 = pc3.localDescription.sdp;
                console.log(answer3)
                answer3 += "&"; 
                ws3.send(answer3);
                ws3.close();
            } else {
                console.log("Retry");
                console.log(pc3.localDescription.sdp);
                restartWebRTC3();
            };
        }
        });
    // pc3.onicecandidate = e3 => {
    //     console.log(e3);
    //     if (e3.candidate == "") {
    //         console.log(answer3);
    //         if (pc3.localDescription.sdp.includes("c=IN IP4 0.0.0.0")) {
    //             answer3 = pc3.localDescription.sdp;
                
    //             answer3 += "&";
    //             ws3.send(answer3);
    //             return;
    //         } else {
    //             console.log("Retry");
    //             createAnswer3();
    //         }
    //     }
    // }
};


function startChannel3() {
    console.log("Trying to start channel")
    if (pc3.signalingState != "have-local-offer") return;
    var desc = new RTCSessionDescription({ type:"answer", sdp:answer3 });
    pc3.setRemoteDescription(desc);
    ws3.close();
};

function sendData3(data) {
    dc3.send(data)
};

function startWebRTC3() {
    var id = document.getElementById("self-id").value;
    if (id == 1) {
        ws3 = new WebSocket("ws://localhost:13374");
        welcomeMsg3 = "Hello from Server 1"
        ws3.onopen = function() {
            createOffer3();
            console.log("WS open");
        };
        ws3.onmessage = function (evt) {
            answer3 = evt.data;
            console.log(answer3);
            startChannel3();
        }
    } else if (id == 4){
        // Wait for offer
        ws3 = new WebSocket("ws://localhost:13374");
        welcomeMsg3 = "Hello from Server 4"
        ws3.onopen = function() {
            ws3.send("setup:recvOffer\n&");
            };
        ws3.onmessage = function (evt) { 
            offer3 = evt.data;
            console.log(offer3);
            createAnswer3();
        };
    }
};

function restartWebRTC3() {
    dc3, pc3 = new RTCPeerConnection({ iceServers: [server] });
    pc3.ondatachannel = e3 => dcInit3(dc3 = e3.channel);
    pc3.oniceconnectionstatechange = e3 => console.log(pc3.iceConnectionState);
        
    offer3 = "";
    answer3 = "";

    startWebRTC3();
};
        