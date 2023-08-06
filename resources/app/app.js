document.querySelector("form").addEventListener("submit", event => {
  event.preventDefault()
  var info = {};

  for (const section of document.querySelectorAll("section")) {
    let input = section.querySelector("input")
    if (!input && section.querySelector("button")) {
      let parentKey = section.getAttribute("data-var-attach").split(":")[0]
      let actualKey = section.getAttribute("data-var-attach").split(":")[1]
      if (!info[parentKey])
        info[parentKey] = {}

      if (section.querySelector("button").textContent === "Bind") {
        info[parentKey][actualKey] = ""
        continue
      }
      info[parentKey][actualKey] = section.querySelector("button").textContent
      continue
    }

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

document.querySelector("#generate-shortcut").addEventListener("click", (event) => {
  event.preventDefault()
  astilectron.sendMessage(JSON.stringify({ "type": "generate-shortcut" }), (m) => {
    console.log(m)
  })
})
document.querySelectorAll("section[data-keyinput]").forEach(elem => {
  const bindButton = elem.querySelector("button")
  bindButton.textContent = "Bind"
  bindButton.addEventListener("click", (event) => {
    event.preventDefault()

    bindButton.textContent = "Press a key..."
    elem.setAttribute("data-keyinput", "true")
  })
})

const modifierUsed = {
  "Shift": false,
  "Control": false,
  "Alt": false,
  "Meta": false
}
document.addEventListener("keydown", (e) => {
  let keyinput = document.querySelector("section[data-keyinput=true]")
  // check if the key is either a shift, ctrl, alt, or meta key
  if (e.key === "Shift" || e.key === "Control" || e.key === "Alt" || e.key === "Meta") {
    modifierUsed[e.key] = true
    return
  }
  if (keyinput) {
    let keybind = ""
    if (modifierUsed["Control"]) keybind += "ctrl+"
    if (modifierUsed["Alt"]) keybind += "alt+"
    if (modifierUsed["Shift"]) keybind += "shift+"
    if (modifierUsed["Meta"]) keybind += "meta+"
    keybind += e.key
    keyinput.querySelector("button").textContent = keybind.toLowerCase()
    keyinput.setAttribute("data-keyinput", "false")
  }
})

document.addEventListener("keyup", (e) => {
  // check if the key is either a shift, ctrl, alt, or meta key
  if (e.key === "Shift" || e.key === "Control" || e.key === "Alt" || e.key === "Meta") {
    modifierUsed[e.key] = false
    return
  }
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
                if (!input && elem.querySelector("button")) {
                  if (valtoset === "") {
                    elem.querySelector("button").textContent = "Bind"
                    return
                  }
                  elem.querySelector("button").textContent = valtoset
                  return
                }
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