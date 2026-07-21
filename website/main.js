const lines = [
    { t: "$ curl -sSL .../install.sh | sh", cls: "prompt" },
    { t: "Downloading pokedexcli-darwin-arm64...", cls: "" },
    { t: "Installed! Run it with: pokedexcli", cls: "" },
    { t: "", cls: "" },
    { t: "$ pokedexcli", cls: "prompt" },
    { t: "PROF. OAK: Welcome to the world of POKEMON!", cls: "" },
    { t: "PROF. OAK: This world is inhabited far and wide", cls: "" },
    { t: "by creatures called POKEMON.", cls: "" }
  ];

  const out = document.getElementById('termOut');
  const reduceMotion = window.matchMedia('prefers-reduced-motion: reduce').matches;

  function renderStatic(){
    out.innerHTML = lines.map(function(l){
      return '<span class="' + l.cls + '">' + l.t + '</span>';
    }).join('\n') + '<span class="cursor"></span>';
  }


  if (reduceMotion) {
    renderStatic();
  } else {
    var li = 0, ci = 0, html = '';
    function typeStep(){
      if (li >= lines.length) {
        out.innerHTML = html + '<span class="cursor"></span>';
        return;
      }
      var line = lines[li];
      if (ci <= line.t.length) {
        var current = html + '<span class="' + line.cls + '">' + line.t.slice(0, ci) + '</span><span class="cursor"></span>';
        out.innerHTML = current;
        ci++;
        setTimeout(typeStep, 14);
      } else {
        html += '<span class="' + line.cls + '">' + line.t + '</span>\n';
        li++; ci = 0;
        setTimeout(typeStep, 90);
      }
    }
    typeStep();
  }
 
  document.querySelectorAll('.copy-btn').forEach(function(btn){
    btn.addEventListener('click', function(){
      var target = document.getElementById(btn.getAttribute('data-target'));
      var text = target.textContent;
      navigator.clipboard.writeText(text).then(function(){
        var original = btn.textContent;
        btn.textContent = 'copied!';
        setTimeout(function(){ btn.textContent = original; }, 1600);
      });
    });
  });


  