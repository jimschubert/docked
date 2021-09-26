FROM scratch

# Recommendation apt
RUN apt install -y \
  bzr \
  git \
  cvs \
  subversion \
  mercurial
