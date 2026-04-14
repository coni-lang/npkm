with open("main.coni", "r") as f: text = f.read()
bad = """                    "echo 'No package manager found' && exit 1")))))]
          res (shell/sh cmd)]"""
good = """                    "echo 'No package manager found' && exit 1"))))))
          res (shell/sh cmd)]"""
if bad in text:
    text = text.replace(bad, good)
    with open("main.coni", "w") as f: f.write(text)
    print("Fixed bracket in PackageTask")
else:
    print("Could not find the target string.")
