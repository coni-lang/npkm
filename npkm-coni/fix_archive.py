with open("main.coni", "r") as f: text = f.read()
bad_zip = """                (str "cd "$(dirname '" (:src s) "')" && zip -r '" (:dest s) "' "$(basename '" (:src s) "')"")"""
good_zip = """                (str "cd \\\"$(dirname '" (:src s) "')\\\" && zip -r '" (:dest s) "' \\\"$(basename '" (:src s) "')\\\"")"""

bad_tar = """                (str "tar -czf '" (:dest s) "' -C "$(dirname '" (:src s) "')" "$(basename '" (:src s) "')""))"""
good_tar = """                (str "tar -czf '" (:dest s) "' -C \\\"$(dirname '" (:src s) "')\\\" \\\"$(basename '" (:src s) "')\\\""))"""

if bad_zip in text and bad_tar in text:
    text = text.replace(bad_zip, good_zip).replace(bad_tar, good_tar)
    with open("main.coni", "w") as f: f.write(text)
    print("Fixed unescaped quotes.")
else:
    print("Could not find the target strings.")
