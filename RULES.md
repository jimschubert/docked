# Rules

* [D7:tagged-latest](#d7:tagged-latest)
* [D7:tagged-latest-builder](#d7:tagged-latest-builder)

## D7:tagged-latest

> _Avoid using images tagged as Latest in production builds_

Docker best practices suggest avoiding &#39;latest&#39; images in production builds

Priority: **High**  
Analyzes: <kbd><a href="https://docs.docker.com/engine/reference/builder/#from">FROM</a></kbd>

## D7:tagged-latest-builder

> _Avoid using images tagged as Latest in builder stages_

Using &#39;latest&#39; images in builders is not recommended.

Priority: **Low**  
Analyzes: <kbd><a href="https://docs.docker.com/engine/reference/builder/#from">FROM</a></kbd>

