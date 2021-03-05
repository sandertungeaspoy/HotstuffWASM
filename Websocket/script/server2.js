var server = { urls: "stun:stun.l.google.com:19302" };

var dc2, pc2 = new RTCPeerConnection({ iceServers: [server] });
pc2.ondatachannel = e2 => dcInit2(dc2 = e2.channel);
pc2.oniceconnectionstatechange = e2 => console.log(pc2.iceConnectionState);


var ws2;

var offer2 = "";
var answer2 = "";
var welcomeMsg2 = "";

function dcInit2() {
    dc2.onopen = () => sendData2(welcomeMsg2);
    dc2.onmessage = e2 => console.log(e2.data);
}

function createOffer2() {
    dcInit2(dc2 = pc2.createDataChannel("chat"));
    pc2.createOffer().then(d2 => pc2.setLocalDescription(d2));
    pc2.onicecandidate = e2 => {
        pc2.addEventListener("icegatheringstatechange", ev => {
            switch(pc2.iceGatheringState) {
              case "new":
                /* gathering is either just starting or has been reset */
                break;
              case "gathering":
                /* gathering has begun or is ongoing */
                break;
              case "complete":
                /* gathering has ended */
                offer2 = pc2.localDescription.sdp;
                // console.log(offer2);
                offer2 += "&"
                ws2.send(offer2);
                // console.log(offer2);
                ws2.send("setup:recvAnswer\n&");
                break;
            }
          });
        // console.log(e2);
        // if (e2.candidate) return;
        // offer2 = pc2.localDescription.sdp;
        // console.log(offer2);
        // offer2 += "&";
        // ws2.send(offer2);
        // console.log(offer2);
        // ws2.send("setup:recvAnswe\n&");
    };
    return true;
};


function createAnswer2() {
    if (pc2.signalingState != "stable") return;
    var desc = new RTCSessionDescription({ type:"offer", sdp:offer2 });
    pc2.setRemoteDescription(desc)
        .then(() => pc2.createAnswer()).then(d2 => pc2.setLocalDescription(d2));

    pc2.addEventListener("icegatheringstatechange", ev => {
        switch(pc2.iceGatheringState) {
            case "new":
            /* gathering is either just starting or has been reset */
            break;
            case "gathering":
            /* gathering has begun or is ongoing */
            break;
            case "complete":
            /* gathering has ended */
            if (pc2.localDescription.sdp.includes("c=IN IP4 0.0.0.0")) {
                answer2 = pc2.localDescription.sdp;
                // console.log(answer2)
                answer2 += "&"; 
                ws2.send(answer2);
                ws2.close();
            } else {
                console.log("Retry");
                // console.log(pc2.localDescription.sdp);
                restartWebRTC2();
            };
        }
    });    
    // pc2.onicecandidate = e2 => {
    //     console.log(e2);
    //     if (e2.candidate = "") {
    //         console.log(answer2);
    //         if (pc2.localDescription.sdp.includes("c=IN IP4 0.0.0.0")) {
    //             answer2 = pc.localDescription.sdp;
                
    //             answer2 += "&";
    //             ws2.send(answer2);
    //             return;
    //         } else {
    //             console.log("Retry");
    //             createAnswer2();
    //         }
    //     }
    // }
};


function startChannel2() {
    console.log("Trying to start channel")
    if (pc2.signalingState != "have-local-offer") return;
    var desc = new RTCSessionDescription({ type:"answer", sdp:answer2 });
    pc2.setRemoteDescription(desc);
    ws2.close();
};

function sendData2(data) {
    dc2.send(data)
};

function startWebRTC2() {
    var id = document.getElementById("self-id").value;
    if (id == 1) {
        ws2 = new WebSocket("ws://localhost:13373");
        welcomeMsg2 = "Hello from Server 1"
        ws2.onopen = function() {
            createOffer2();
            console.log("WS open");
            };
        ws2.onmessage = function (evt) {
            console.log("msg recv");
            // console.log(evt);
            answer2 = evt.data;
            // console.log(answer2);
            startChannel2();
        }
    } else if (id == 3) {
        // Wait for offer
        welcomeMsg2 = "Hello from Server 3"
        ws2 = new WebSocket("ws://localhost:13373");
        ws2.onopen = function() {
            ws2.send("setup:recvOffer\n&");
            };
        ws2.onmessage = function (evt) { 
            offer2 = evt.data;
            // console.log(offer2);
            createAnswer2();
        };
    }
};

function restartWebRTC2() {
    dc2, pc2 = new RTCPeerConnection({ iceServers: [server] });
    pc2.ondatachannel = e2 => dcInit2(dc2 = e2.channel);
    pc2.oniceconnectionstatechange = e2 => console.log(pc2.iceConnectionState);
        
    offer2 = "";
    answer2 = "";

    startWebRTC2();
};
    