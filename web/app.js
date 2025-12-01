const form = document.getElementById("character-form");
const statusEl = document.getElementById("status");
const resultCard = document.getElementById("result-card");
const resultBasic = document.getElementById("result-basic");
const resultAbilities = document.getElementById("result-abilities");
const resultCombat = document.getElementById("result-combat");
const resultId = document.getElementById("result-id");

form.addEventListener("submit", async (e) => {
  e.preventDefault();

  statusEl.textContent = "Генерирую персонажа...";
  statusEl.className = "status";

  const submitBtn = form.querySelector("button[type=submit]");
  submitBtn.disabled = true;

  const data = {
    name: form.name.value.trim(),
    class: form.class.value.trim(),
    race: form.race.value.trim(),
    background: form.background.value.trim(),
    level: Number(form.level.value) || 1,
  };

  try {
    const res = await fetch("/characters/generate", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(data),
    });

    const payload = await res.json().catch(() => ({}));

    if (!res.ok) {
      const msg = payload.error || "Не удалось создать персонажа";
      statusEl.textContent = msg;
      statusEl.className = "status error";
      resultCard.hidden = true;
      return;
    }

    renderResult(payload);
    statusEl.textContent = "Персонаж создан";
    statusEl.className = "status success";
  } catch (err) {
    console.error(err);
    statusEl.textContent = "Ошибка сети";
    statusEl.className = "status error";
    resultCard.hidden = true;
  } finally {
    submitBtn.disabled = false;
  }
});

function renderResult(c) {
  resultCard.hidden = false;

  resultBasic.innerHTML = "";
  resultBasic.append(
    pill(`Имя: ${escapeHtml(c.name || "")}`),
    pill(`Класс: ${escapeHtml(c.class || "")}`),
    pill(`Раса: ${escapeHtml(c.race || "")}`),
    pill(`Предыстория: ${escapeHtml(c.background || "")}`),
    pill(`Уровень: ${c.level ?? "?"}`),
  );

  resultAbilities.innerHTML = "";
  if (c.abilityScores) {
    const a = c.abilityScores;
    addAbility("STR", a.strength);
    addAbility("DEX", a.dexterity);
    addAbility("CON", a.constitution);
    addAbility("INT", a.intelligence);
    addAbility("WIS", a.wisdom);
    addAbility("CHA", a.charisma);
  }

  resultCombat.innerHTML = "";
  addStat("AC", c.armorClass);
  addStat("HP", `${c.currentHitPoints}/${c.maxHitPoints}`);
  addStat("Temp HP", c.temporaryHitPoints);
  addStat("Speed", `${c.speed} ft`);
  addStat("Prof", formatSigned(c.proficiencyBonus));
  addStat("Init", formatSigned(c.initiative));

  resultId.textContent = c.id || "—";
}

function pill(text) {
  const span = document.createElement("span");
  span.className = "pill";
  span.textContent = text;
  return span;
}

function addAbility(label, score) {
  const wrap = document.createElement("div");
  wrap.className = "ability";

  const lab = document.createElement("div");
  lab.className = "ability-label";
  lab.textContent = label;

  const val = document.createElement("div");
  val.className = "ability-value";
  val.textContent = score ?? "—";

  const mod = document.createElement("div");
  mod.className = "ability-mod";
  mod.textContent = score ? formatSigned(abilityModifier(score)) : "";

  wrap.append(lab, val, mod);
  resultAbilities.appendChild(wrap);
}

function addStat(label, value) {
  const wrap = document.createElement("div");
  wrap.className = "stat";

  const lab = document.createElement("div");
  lab.className = "stat-label";
  lab.textContent = label;

  const val = document.createElement("div");
  val.className = "stat-value";
  val.textContent = value ?? "—";

  wrap.append(lab, val);
  resultCombat.appendChild(wrap);
}

function abilityModifier(score) {
  return Math.floor((score - 10) / 2);
}

function formatSigned(v) {
  if (typeof v !== "number" || Number.isNaN(v)) return "—";
  return v >= 0 ? `+${v}` : `${v}`;
}

function escapeHtml(str) {
  return String(str)
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&#039;");
}


