# Updating github's public key
ssh-keygen -R github.com
ssh-keyscan github.com >> ~/.ssh/known_hosts

# Required parameters:
#   PROJECT_OWNER - Project's owner username
#   PROJECT_NAME - Project's repository name
#   GOLANG_VERSION - Project's golang version
#   SECRETS - Project's secrets
# Additional parameters
#   PROJECT_SECRETS_NAME - Project's secrets filename
#   BRANCH_NAME - Branch name to update from

[[ -n $PROJECT_NAME ]] || exit 1
[[ -n $PROJECT_OWNER ]] || exit 1
[[ -n $GOLANG_VERSION ]] || exit 1

PROJECT_DIR=/srv/$PROJECT_NAME
PROJECT_LOG_BASE_DIR=/var/log/$PROJECT_NAME

APPEND_DEPENDENCIES=(
  bison
  git
)
APPEND_LOG_DIRS=(
  "systemd/echochamber-dc"
)
PROJECT_SECRETS_NAME='.env'

[[ -n $BRANCH_NAME ]] || BRANCH_NAME='main'

# Updating dependencies
echo -e "Install apts"
apt update --fix-missing
apt -y install "${APPEND_DEPENDENCIES[@]}"

echo -e "Deploy ${PROJECT_NAME} from branch ${BRANCH_NAME}"

# Acquiring repository data
echo -e "Clone repository"
cd /srv || exit 1
echo -e "Using public link to clone repository"
GIT_REPOSITORY="https://github.com/${PROJECT_OWNER}/${PROJECT_NAME}.git"

if ! git clone --branch $BRANCH_NAME "$GIT_REPOSITORY"; then
  cd "$PROJECT_NAME" || exit 1
  git pull origin $BRANCH_NAME
else
  git config --global --add safe.directory "$PROJECT_DIR"
fi

# Installing GVM
if ! command -v gvm &>/dev/null 2>&1; then
  echo -e "Installing GVM"
  bash < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)
  source ~/.gvm/scripts/gvm
fi

# Installing proper golang version
echo -e "Setting up go$GOLANG_VERSION"
gvm install go$GOLANG_VERSION
gvm use go$GOLANG_VERSION

# Including utils
. "$PROJECT_DIR"/deploy/deploy_utils.sh

# Updating secrets
echo -e "Set secrets"
echo "$SECRETS" > "$PROJECT_DIR"/$PROJECT_SECRETS_NAME

# Creating user to run the project with
echo -e "Create user"
useradd --no-create-home -d /nonexistent "$PROJECT_NAME"

# Building binary
echo -e "Build binary"
cd $PROJECT_DIR
rm -f "$PROJECT_DIR/binary"
go build -ldflags="-s -w" -o "$PROJECT_DIR/binary" "$PROJECT_DIR/src/"
chmod +x "$PROJECT_DIR/binary"

# Creating directories
echo -e "Create dirs"
create_log_dirs "$PROJECT_LOG_BASE_DIR" "${APPEND_LOG_DIRS[@]}"

# Setup directory permissions
echo -e "Setup permissions"
chown -R "$PROJECT_NAME":"$PROJECT_NAME" "$PROJECT_DIR" "$PROJECT_LOG_BASE_DIR"

PROJECT_DEPLOY_PATH=$PROJECT_DIR/deploy

SYSTEMD_PATH=/etc/systemd/system

# Applying execution permissions
echo -e "Applying execution permissions"
apply_execution_permissions "$PROJECT_DEPLOY_PATH" "${APPEND_EXECUTION_PERMISSION_DIRS[@]}"

# Update systemd demons & run it
echo -e "Update systemd ${PROJECT_DEPLOY_SECOND_PATH}"

update_symlinks "$PROJECT_DEPLOY_PATH"/systemd $SYSTEMD_PATH

systemctl enable "$PROJECT_NAME".target
systemctl daemon-reload
systemctl restart "$PROJECT_NAME".target

# Update logrotate settings
echo -e "Update logrotate"
LOG_ROTATE_PATH=/etc/logrotate.d/
update_symlinks "$PROJECT_DEPLOY_PATH"/logrotate $LOG_ROTATE_PATH true

exit 0
