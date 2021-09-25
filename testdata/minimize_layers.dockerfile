# minimize-layers will emit a recommendation when the combination of all
# RUN, COPY, ADD layers is greater than 6. This allows for 2 of each, but
# the number of each is not enforced. You could have 3 RUN and 4 COPY, or
# just 7 RUN to trigger the recommendation.
FROM scratch

RUN mkdir /a
RUN mkdir /b
RUN mkdir /c
COPY LICENSE /a/LICENSE
COPY README.md /b/README.md
ADD TODO.md /b/TODO.md
ADD Dockerfile /a/Dockerfile
