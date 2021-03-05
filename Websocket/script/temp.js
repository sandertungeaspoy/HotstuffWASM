var server = { urls: "stun:stun.l.google.com:19302" };

var dc, pc = new RTCPeerConnection({ iceServers: [server] });
pc.ondatachannel = e => dcInit(dc = e.channel);
pc.oniceconnectionstatechange = e => log(pc.iceConnectionState);


function dcInit() {
  dc.onopen = () => log("Chat!");
  dc.onmessage = e => log(e.data);
}

function createOffer() {
  dcInit(dc = pc.createDataChannel("chat"));
  pc.createOffer().then(d => pc.setLocalDescription(d));
  pc.onicecandidate = e => {
    if (e.candidate) return;
    return pc.localDescription.sdp;
  };
};


function createAnswer(offer) {
  if (pc.signalingState != "stable") return;
  var desc = new RTCSessionDescription({ type:"offer", offer });
  pc.setRemoteDescription(desc)
    .then(() => pc.createAnswer()).then(d => pc.setLocalDescription(d));
  pc.onicecandidate = e => {
    if (e.candidate) return;
    return pc.localDescription.sdp;
  };
};


function startChannel(answer) {
  if (pc.signalingState != "have-local-offer") return;
  answer.disabled = true;
  var desc = new RTCSessionDescription({ type:"answer", answer });
  pc.setRemoteDescription(desc);
};

function sendData(data) {
  dc.send(data)
};



