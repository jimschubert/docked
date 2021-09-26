FROM scratch

# Success apt-get
RUN apt-get install -y \
  bzr \
  cvs \
  git \
  mercurial \
  subversion
