import re

with open("/Users/nico/cool/npkm/npkm-coni/main.coni", "r") as f:
    content = f.read()

# 1. Add slice-plays-by-task-idx
if "defn slice-plays-by-task-idx" not in content:
    slice_func = """(defn slice-plays-by-task-idx [plays start-idx]
  (if (<= start-idx 0)
    plays
    (let [state (atom 0)]
      (map (fn [play]
             (let [filtered-tasks (filter (fn [t]
                                            (let [curr @state]
                                              (swap! state (fn [x] (+ x 1)))
                                              (>= curr start-idx)))
                                          (:tasks play))]
               (assoc play :tasks filtered-tasks)))
           plays))))
"""
    content = content.replace("(defn execute-playbook ", slice_func + "\n(defn execute-playbook ")

# 2. Update routing base path and start-idx
routing_old = """              (let [path (get req :path)
                    method (get req :method)]
                (if (= path "/")"""
routing_new = """              (let [raw-path (get req :path)
                    path-parts (str/split raw-path "?")
                    path (first path-parts)
                    q-str (if (> (count path-parts) 1) (second path-parts) "")
                    start-idx (if (str/starts-with? q-str "start-idx=") (sys-parse-int (second (str/split q-str "="))) 0)
                    method (get req :method)]
                (if (= path "/")"""
if routing_old in content:
    content = content.replace(routing_old, routing_new)

# 3. Update the execute-playbook call in /api/run to use slice-plays-by-task-idx
api_run_old = """                            (let [content (io/read-file playbook-file)
                                  parsed-data (parse-playbook playbook-file content)
                                  tasks (:tasks parsed-data)
                                  cfg (:cfg parsed-data)
                                  inventory nil]
                              (execute-playbook tasks inventory cfg false content true false false))"""
api_run_new = """                            (let [content (io/read-file playbook-file)
                                  parsed-data (parse-playbook playbook-file content)
                                  tasks (if (map? parsed-data)
                                          (:tasks parsed-data)
                                          parsed-data)
                                  sliced-tasks (slice-plays-by-task-idx tasks start-idx)
                                  cfg (:cfg parsed-data)
                                  inventory nil]
                              (execute-playbook sliced-tasks inventory cfg false content true false false))"""
if api_run_old in content:
    content = content.replace(api_run_old, api_run_new)

# 4. Patch HTML
html_match = re.search(r':body \(str "(<!DOCTYPE html>.*?)(\"\)})', content, flags=re.DOTALL)
if html_match:
    html = html_match.group(1)
    
    # Update Start button
    html = html.replace('id=\\"start-btn\\" class=\\"btn\\" style=\\"width: 100%; margin-bottom: 10px;\\" onclick=\\"startRun()\\">Start Execution',
                        'id=\\"start-btn\\" class=\\"btn\\" style=\\"width: 100%; margin-bottom: 10px;\\" onclick=\\"configureExecution()\\">Execution')
    
    # Add Arch Map button
    if 'id=\\"arch-btn\\"' not in html:
        html = html.replace('onclick=\\"runAudit()\\">Run Security Audit</button>',
                            'onclick=\\"runAudit()\\">Run Security Audit</button><button id=\\"arch-btn\\" class=\\"btn\\" style=\\"width: 100%; margin-bottom: 10px; background: #f59e0b; box-shadow: 0 4px 15px rgba(245, 158, 11, 0.4);\\" onclick=\\"viewArchitecture()\\">Architecture Map</button>')
    
    # Add JS functions
    if "function configureExecution" not in html:
        js_to_add = r"""function configureExecution() { document.getElementById('start-btn').disabled = false; document.getElementById('audit-btn').disabled = false; document.getElementById('logs-container').innerHTML = '<p style=\"color:#64748b; text-align:center; margin-top:100px;\">Loading Configuration...</p>'; fetch('/api/info').then(r => r.json()).then(data => { let html = '<h2 style=\"margin-top:0; border:none; font-size:2rem; font-weight:800; color:transparent; background:linear-gradient(to right, #38bdf8, #818cf8); -webkit-background-clip:text;\">Execution Configuration</h2><div class=\"task-card\" style=\"padding: 25px;\"><label style=\"color:#cbd5e1; font-weight:600; display:block; margin-bottom:10px;\">Start at Task:</label><select id=\"start-task-select\" style=\"width:100%; padding:10px; background:rgba(0,0,0,0.5); border:1px solid rgba(255,255,255,0.1); color:#f8fafc; border-radius:6px; margin-bottom:20px;\"><option value=\"0\">-- Beginning of Playbook --</option>'; if (data.tasks) { let plays = Array.isArray(data.tasks) ? data.tasks : [data.tasks]; let absIdx = 0; plays.forEach((play, pIdx) => { if (play.tasks) { play.tasks.forEach((t, tIdx) => { html += `<option value=\"${absIdx}\">[Play: ${play.name || pIdx}] ${t.name || 'Unnamed Task'}</option>`; absIdx++; }); } }); } html += '</select><button class=\"btn\" style=\"width:100%; background:#4ade80; color:#064e3b;\" onclick=\"startRun(document.getElementById(\\'start-task-select\\').value)\">Run Playbook</button></div>'; document.getElementById('logs-container').innerHTML = html; }); } function viewArchitecture() { document.getElementById('logs-container').innerHTML = '<p style=\"color:#64748b; text-align:center; margin-top:100px;\">Loading Architecture Map...</p>'; fetch('/api/audit').then(r => r.json()).then(data => { document.getElementById('logs-container').innerHTML = '<div style=\"background: rgba(0,0,0,0.3); border-radius: 8px; padding: 20px; width: 100%; min-height: 500px; display: flex; justify-content: center; align-items: center;\"><div class=\"mermaid\">graph TD\\\\n' + data.mermaid + '</div></div>'; if(window.mermaid) { window.mermaid.run(); } }); } """
        html = html.replace('function startRun() {', js_to_add + 'function startRun(idx = 0) {')
        html = html.replace("const es = new EventSource('/api/run');", "const es = new EventSource('/api/run?start-idx=' + idx);")
        if "document.getElementById('arch-btn')" not in html:
            html = html.replace("document.getElementById('browse-btn').disabled = true;", "document.getElementById('browse-btn').disabled = true; document.getElementById('arch-btn').disabled = true;")
            html = html.replace("document.getElementById('browse-btn').disabled = false;", "document.getElementById('browse-btn').disabled = false; document.getElementById('arch-btn').disabled = false;")

    content = content[:html_match.start(1)] + html + content[html_match.end(1):]
else:
    print("HTML not found!")

with open("/Users/nico/cool/npkm/npkm-coni/main.coni", "w") as f:
    f.write(content)

print("Patch applied successfully!")
