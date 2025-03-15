window.onload = ()=>{

    let commandBox = document.getElementById("command-input-box")
    let caret = document.getElementById("caret-box");
    let history = document.getElementById("history");
    function getCaretCoordinates(){
        const selection = window.getSelection();
        if(!selection.rangeCount) return getFallbackCoords();

        const range = selection.getRangeAt(0);
        const rect = range.getBoundingClientRect();

        if (rect.width === 0 && rect.height === 0) {
            return getFallbackCoords();
        }

        return {x: rect.left + window.scrollX, y: rect.top + window.scrollX}
    }

    function getFallbackCoords() {
        const commandBoxCoords = commandBox.getBoundingClientRect();
        return { x: commandBoxCoords.left + window.scrollX, y: commandBoxCoords.top + window.scrollY };
    }

    function RepositionCaret(e){
        setTimeout(()=>{
            const {x, y} = getCaretCoordinates();
            caret.style.top = `${y}px`;
            caret.style.left = `${x}px`;
        }, 0);
    }

    function sendCommand(){
        let command = commandBox.innerText;
        commandBox.innerText = ""
        let result = fetch('/command', {
            method: "POST",
            headers: {
                "Content-Type": "text/plain"
            },
            body: command
        }).then((res)=>{
            return res.text();
        }).then((res)=>{
            history.innerHTML += res + "<br>"
        }).catch((err)=>{
            console.log('caught', err)
        })
    }

    commandBox.addEventListener("keydown", (e)=>{
        if(e.key==="Enter"){
            e.preventDefault()
            sendCommand()
        }
        RepositionCaret()
    })
    commandBox.addEventListener("click", RepositionCaret)
    RepositionCaret()
}