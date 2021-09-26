FROM scratch

# Recommendation apt-get
RUN apt-get install -y \
  bzr \
  git \
  cvs \
  subversion \
  mercurial
