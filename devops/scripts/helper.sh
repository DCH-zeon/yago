#!/bin/bash

# Color codes
RED='\033[0;31m'
GREEN='\033[1;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;36m'
NC='\033[0m' # No Color
S='\033[38;5;63m'
G1='\033[38;5;45m'
G2='\033[38;5;44m'
G3='\033[38;5;43m'
G4='\033[38;5;42m'
G5='\033[38;5;41m'
G6='\033[38;5;40m'
G7='\033[38;5;02m'

yago() {
  echo -e "${G1}   ██${S}╗${G1}   ██${S}╗${G1}  █████${S}╗${G1}   ██████${S}╗${G1}   ██████${S}╗${NC}"
  echo -e "${G2}   ██${S}║${G2}   ██${S}║${G2} ██${S}╔══${G2}██${S}╗${G2} ██${S}╔═══${G2}██${S}╗${G2} ██${S}╔═══${G2}██${S}╗${NC}"
  echo -e "${G3}   ${S}╚${G3}██${S}╗${G3} ██${S}╔╝${G3} ██${S}║${G3}  ██${S}║${G3} ██${S}║${G3}   ${S}╚═╝${G3} ██${S}║${G3}   ██${S}║${NC}"
  echo -e "${G4}    ${S}╚${G4}████${S}╔╝${G4}  ███████${S}║${G4} ██${S}║${G4}  ███${S}╗${G4} ██${S}║${G4}   ██${S}║${NC}"
  echo -e "${G5}     ${S}╚${G5}██${S}╔╝${G5}   ██${S}╔══${G5}██${S}║${G5} ██${S}║${G5}   ██${S}║${G5} ██${S}║${G5}   ██${S}║${NC}"
  echo -e "${G6}    ██████${S}╗${G6}  ██${S}║${G6}  ██${S}║ ╚${G6}███████${S}║ ╚${G6}██████${S}╔╝${NC}"
  echo -e "${S}    ╚═════╝  ╚═╝  ╚═╝  ╚══════╝   ╚════╝${NC}"
  echo -e "${G7}   ======================================"
}

cat_watch() {
  frames=(
  "( -.- )"
  "( o.o )"
  "( o.o )"
  "( -.- )"
  "( o.o )"
  "( o.o )"
  "( o.o )"
  )
  echo " /\\_/\\"
  for i in {1..10}; do
    for f in "${frames[@]}"; do
      echo -ne "\r$f\n"
      echo -ne "▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔ \033[A"
      sleep 0.2
    done
    sleep 0.3
  done
  echo ""
  echo ""
  echo "cat is watching you"
}

generate_ssl() {
  local CERT_DIR="./devops/nginx/ssl"
  local DOMAIN="${DOMAIN:-yago.loc}"
  local IS_WSL=0
  local MKCERT_BIN="mkcert"

  if grep -qi microsoft /proc/version 2>/dev/null; then
    IS_WSL=1
    MKCERT_BIN="./devops/nginx/mkcert.exe"
  fi
  echo "🔐 Installing a local CA via mkcert..."
  if [ "$IS_WSL" -eq 1 ]; then
    if [ ! -f "$MKCERT_BIN" ]; then
      echo -e "${RED}❌ mkcert.exe not found: $MKCERT_BIN${NC}"
      exit 1
    fi
    powershell.exe -ExecutionPolicy Bypass -File "$(wslpath -w "$(pwd)/devops/scripts/generate-ssl.ps1")" -Domain "${DOMAIN:-yago.loc}" -Wait >/dev/null 2>&1
  else
    if ! command -v mkcert >/dev/null 2>&1; then
      if ! "$MKCERT_BIN" -install; then
        echo -e "${RED}❌ mkcert -install failed${NC}"
        echo -e "${YELLOW}⚠️ Please install mkcert manually${NC}"
        exit 1
      fi
      if ! "$MKCERT_BIN" \
        -cert-file "$CERT_DIR/$DOMAIN.crt" \
        -key-file "$CERT_DIR/$DOMAIN.key" \
        "$DOMAIN" \
        "*.$DOMAIN" \
        "admin.$DOMAIN" \
        "api.$DOMAIN"; then
        echo -e "${RED}❌ Failed to generate certificates${NC}"
        exit 1
      fi
    fi
  fi

  if [ "$IS_WSL" -eq 1 ]; then
    WIN_HOSTS="/mnt/c/Windows/System32/drivers/etc/hosts"
    if ! grep -q "$DOMAIN" "$WIN_HOSTS"; then
        WIN_PATH=$(wslpath -w "$(pwd)/devops/scripts/add-hosts.ps1")
        powershell.exe -Command "Start-Process PowerShell -Verb RunAs -WindowStyle Hidden -ArgumentList '-NoProfile -ExecutionPolicy Bypass -File \"$WIN_PATH\" -Domain $DOMAIN' -Wait"
        WIN_CERT=$(wslpath -w "$CERT_DIR/$DOMAIN.crt")
        powershell.exe -Command "Start-Process PowerShell -Verb RunAs -WindowStyle Hidden -ArgumentList '-NoProfile Import-Certificate -FilePath \"$WIN_CERT\" -CertStoreLocation Cert:\\LocalMachine\\Root' -Wait"
    fi
    if grep -q "$DOMAIN" "$WIN_HOSTS" && grep -q "admin.$DOMAIN" "$WIN_HOSTS" && grep -q "api.$DOMAIN" "$WIN_HOSTS"; then
        echo -e "${GREEN}✅ Domain $DOMAIN and all subdomains already in Windows hosts${NC}"
    else
        echo -e "${YELLOW}⚠️ Please check your hosts file for entries:${NC}"
        echo "        127.0.0.1    $DOMAIN"
        echo "        127.0.0.1    admin.$DOMAIN"
        echo "        127.0.0.1    api.$DOMAIN"
        echo -e "⚠️ If they are missing, add them manually an path: ${BLUE}Windows/System32/drivers/etc/hosts${NC}"
        exit 1
    fi
  else
    if ! grep -q "$DOMAIN" /etc/hosts 2>/dev/null; then
      echo "127.0.0.1    $DOMAIN" | sudo tee -a /etc/hosts > /dev/null
      echo "127.0.0.1    admin.$DOMAIN" | sudo tee -a /etc/hosts > /dev/null
      echo "127.0.0.1    api.$DOMAIN" | sudo tee -a /etc/hosts > /dev/null
      echo -e "${GREEN}✅ Hosts entries added${NC}"
    else
      echo -e "${GREEN}✅ Domain $DOMAIN already in hosts${NC}"
    fi
  fi

  echo -e "${GREEN}✅ SSL certificates generated done${NC}"
}

cleanup_ui() {
    if [ -t 1 ]; then
        tput cnorm 2>/dev/null || true
    fi
    printf "\n"
}

is_tty() {
    [ -t 1 ]
}

hide_cursor() {
    if is_tty; then
        tput civis 2>/dev/null || true
    fi
}

show_cursor() {
    if is_tty; then
        tput cnorm 2>/dev/null || true
    fi
}

repeat_char() {
    local char="$1"
    local count="$2"
    local i
    for ((i = 0; i < count; i++)); do
        printf "%s" "$char"
    done
}

draw_bar() {
    local current="$1"
    local total="$2"
    local width=28

    local filled=$(( current * width / total ))
    local empty=$(( width - filled ))

    printf "["
    repeat_char "█" "$filled"
    repeat_char "░" "$empty"
    printf "]"
}

render_frame_once() {
    printf "    /\\_/\\ \n"
    printf "   ( o.o ) \n"
    printf "╔══════════════════════════════════════════════╗\n"
    printf "║       📦 Starting Docker containers          ║\n"
    printf "╟──────────────────────────────────────────────╢\n"
    printf "║                                              ║\n"
    printf "║                                              ║\n"
    printf "║                                              ║\n"
    printf "║                                              ║\n"
    printf "╚══════════════════════════════════════════════╝\n"
}

render_status_block() {
    local stage="$1"
    local frame="$2"
    local elapsed="$3"
    local step="$4"
    local total="$5"
    local frame_cat="$6"

    local percent=$(( step * 100 / total ))
    local bar
    bar=$(draw_bar "$step" "$total")

    # Возврат к началу блока статуса сразу после рамки и заголовка
    tput rc 2>/dev/null || printf "\033[u"
    printf "\n"
    printf "   %s    \n" "$frame_cat"
    printf "\n\n\n"
    # Первая строка контента
    printf "║   %-2s %-35.35s      ║\n" "$frame" "$stage"
    # Вторая строка
    printf "║   progress: %3d%%   elapsed: %4ss            ║\n" "$percent" "$elapsed"
    # Третья строка
    printf "║   %s             ║\n" "$bar"
    # Четвертая строка
    printf "║   logs: %-34.34s   ║\n" "$LOG_FILE"
}

animated_wait() {
    local pid="$1"
    shift
    local stages=("$@")

    local frames=("⢀" "⣀" "⣄" "⣤" "⣴" "⣶" "⣷" "⣿" "⡿" "⡟" "⠟" "⠛" "⠙" "⠉")
    local frame_index=0
    local stage_index=0
    local elapsed=0
    local total=${#stages[@]}
    local frames_cat=("( -.- )" "( o.o )" "( o.o )" "( o.o )" "( o.o )" "( o.o )" "( o.o )")
    hide_cursor

    if is_tty; then
      tput sc 2>/dev/null || printf "\033[s"
      render_frame_once
    fi

    while kill -0 "$pid" 2>/dev/null; do
        local frame="${frames[$((frame_index % ${#frames[@]}))]}"
        local frame_cat="${frames_cat[$((frame_index % ${#frames_cat[@]}))]}"
        local stage="${stages[$stage_index]}"

        if is_tty; then
            render_status_block "$stage" "$frame" "$elapsed" "$((stage_index + 1))" "$total" "$frame_cat"
        else
            printf "\r⠋ %s..." "$stage"
        fi

        frame_index=$((frame_index + 1))
        stage_index=$(( (stage_index + 1) % total ))
        elapsed=$((elapsed + 1))
        sleep 0.12
    done

    show_cursor
}

print_success_screen() {
    if is_tty; then
        tput rc 2>/dev/null || printf "\033[u"
        printf "    /\\_/\\  \n"
        printf "   ( o.o ) \n"
        echo -e "╭──────────────────────────────────────────────╮"
        echo -e "│   ✅ ${GREEN}Containers started successfully${NC}         │"
        echo -e "│   💽 DB:       ${GREEN}ready${NC}                         │"
        echo -e "│   ⚙️ Admin:    ${GREEN}ready${NC}                         │"
        echo -e "│   📍 API:      ${GREEN}ready${NC}                         │"
        echo -e "│   🧩 Nginx:    ${GREEN}ready${NC}                         │"
        echo -e "│   🐇 RabbitMQ: ${GREEN}ready${NC}                         │"
        echo -e "╰──────────────────────────────────────────────╯"
    else
        echo "✅ Containers started successfully"
    fi
}

print_error_screen() {
    echo ""
    echo -e "${RED}❌ Failed to start containers${NC}"
    echo ""
    echo "Check logs with: docker compose logs"
    echo "Or inspect build output: $LOG_FILE"
}

run_with_animation() {
    local cmd=("$@")

    : >"$LOG_FILE"

    if is_tty; then
        "${cmd[@]}" >"$LOG_FILE" 2>&1 &
        local pid=$!

        animated_wait "$pid" \
            "Building images" \
            "Creating containers" \
            "Starting services" \
            "Waiting for readiness"

        wait "$pid"
        return $?
    else
        echo "⏳ Building images..."
        "${cmd[@]}" >"$LOG_FILE" 2>&1
    fi
}
