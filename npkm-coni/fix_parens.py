with open("main.coni", "r") as f: text = f.read()
bad = """                    "echo 'No package manager found' && exit 1"))))))"""
good = """                    "echo 'No package manager found' && exit 1")))))"""
if bad in text:
    text = text.replace(bad, good)
    with open("main.coni", "w") as f: f.write(text)
    print("Fixed parens.")
else: print("Target not found.")
