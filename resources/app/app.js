document.querySelector("form").addEventListener("submit", event => {
  event.preventDefault()
  var info = {};

  for (const section of document.querySelectorAll("section")) {
    let input = section.querySelector("input")
    let parentKey = input.name.split(":")[0]
    let actualKey = input.name.split(":")[1]
    if (!info[parentKey])
      info[parentKey] = {}
    if (section.hasAttribute("data-toggle")) {
      info[parentKey][actualKey] = input.checked
    } else if (section.hasAttribute("data-slider")) {
      info[parentKey][actualKey] = input.value * (section.getAttribute("data-convdecimal") === "%" ? .01 : 1)
    } else {
      info[parentKey][actualKey] = input.value
    }
  }
  console.log({ "type": "saveclose", "value": info })
  astilectron.sendMessage(JSON.stringify({ "type": "saveclose", "value": info }), (m) => {
    if (m === "SUCCESS") window.close()
  })
})
document.querySelectorAll("input[type=range]").forEach(elem => {
  elem.addEventListener("input", () => {
    elem.parentElement.querySelector("label").textContent = elem.value + elem.getAttribute("data-suffix")
  })
})

document.addEventListener('astilectron-ready', function(e) {
  console.log("astilectron-ready")
  // This will listen to messages sent by GO
  astilectron.onMessage(function(message_s) {
    try {
      let message = JSON.parse(message_s)
      if (message.type.startsWith("load")) {
        if (!message.type.endsWith("FIRSTLOAD"))
          document.querySelectorAll("section").forEach(elem => {
            let tagspl = elem.getAttribute("data-var-attach").split(":")
            if (tagspl) {
              let field = tagspl[1]
              let category = tagspl[0]
              let valtoset = message.value[category][field]
              if (valtoset !== undefined) {
                let input = elem.querySelector("input")
                input.setAttribute("name", elem.getAttribute("data-var-attach"))
                if (elem.hasAttribute("data-toggle")) {
                  input.checked = valtoset
                } else if (elem.hasAttribute("data-slider") || elem.hasAttribute("data-textinput")) {
                  if (elem.getAttribute("data-convdecimal") === "%" && ("" + valtoset).includes(".") && !isNaN(valtoset)) {
                    valtoset *= 100
                  }
                  input.value = valtoset
                  if (elem.hasAttribute("data-slider")) {
                    let label = elem.querySelector("label")
                    label.textContent = valtoset + input.getAttribute("data-suffix")
                  }
                }
              }
            }
          })
        return "Loaded Config to UI."
      }
    } catch (error) {
      console.error(error)
      return "" + error
    }

  });
})