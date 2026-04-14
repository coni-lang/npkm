with open("main.coni", "r") as f: text = f.read()
bad = """(defn parse-playbook [file content]
  (let [local-cfg (extract-config content)
        ext-cfg (if (io/exists? "config.yml") (extract-config (io/read-file "config.yml")) {})
        cfg (merge ext-cfg local-cfg)
        interp-content (interpolate-config content cfg)]"""
good = """(defn parse-playbook [file content]
  (let [local-cfg (extract-config content)
        ext-content (if (io/exists? "config.yml") (io/read-file "config.yml") "")
        ext-cfg (if (> (count ext-content) 0) (extract-config ext-content) {})
        cfg (loop [k-list (keys local-cfg) acc ext-cfg]
              (if (empty? k-list) acc
                (recur (rest k-list) (assoc acc (first k-list) (get local-cfg (first k-list))))))
        interp-content (interpolate-config content cfg)]"""

text = text.replace(bad, good)
with open("main.coni", "w") as f: f.write(text)
