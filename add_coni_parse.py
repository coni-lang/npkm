import re

with open("npkm-coni/main.coni", "r") as f:
    text = f.read()

config_code = """
(defn extract-config [content]
  (let [lines (str/split content "\\n")]
    (loop [rem lines
           in-config false
           cfg {}]
      (if (empty? rem)
        cfg
        (let [line (first rem)
              trim-line (str/trim line)]
          (if (= trim-line "config:")
            (recur (rest rem) true cfg)
            (if (or (= trim-line "tasks:") (str/starts-with? trim-line "- name:"))
              (recur (rest rem) false cfg)
              (if (and in-config (str/includes? trim-line ":"))
                (let [colon-idx (str/index-of trim-line ":")
                      k-str (str/trim (str/substring trim-line 0 colon-idx))
                      v-str (str/trim (str/substring trim-line (+ colon-idx 1) (count trim-line)))
                      v-clean (if (and (str/starts-with? v-str "\\\"") (str/ends-with? v-str "\\\""))
                                  (str/substring v-str 1 (- (count v-str) 1))
                                  (if (and (str/starts-with? v-str "'") (str/ends-with? v-str "'"))
                                      (str/substring v-str 1 (- (count v-str) 1))
                                      v-str))]
                  (recur (rest rem) true (assoc cfg k-str v-clean)))
                (recur (rest rem) in-config cfg)))))))))

(defn interpolate-config [content cfg]
  (let [k-list (keys cfg)]
    (loop [rem-keys k-list
           curr content]
      (if (empty? rem-keys)
        curr
        (let [k (first rem-keys)
              v (get cfg k)
              placeholder (str "config." k)
              next-curr (str/replace curr placeholder v)]
          (recur (rest rem-keys) next-curr))))))

(defn parse-playbook [file content]
  (let [local-cfg (extract-config content)
        ext-cfg (if (io/exists? "config.yml") (extract-config (io/read-file "config.yml")) {})
        cfg (merge ext-cfg local-cfg)
        interp-content (interpolate-config content cfg)]
    (if (or (str/ends-with? file ".yml") (str/ends-with? file ".yaml"))
      (read-string (yaml-to-edn interp-content))
      (read-string interp-content))))
"""

orig = """(defn parse-playbook [file content]
  (if (or (str/ends-with? file ".yml") (str/ends-with? file ".yaml"))
    (read-string (yaml-to-edn content))
    (read-string content)))"""

text = text.replace(orig, config_code)

with open("npkm-coni/main.coni", "w") as f:
    f.write(text)

