CSS Optimizer
=============

This is a package (library) and CLI tool for optimizing your CSS


Goals
-----

* Optimizations
	* [ ] Merge `@import`ed CSS files adding comments denoting such
	* [ ] Compile CSS Variables
* Minification Options
	* [ ] Remove comments
		* [ ] Leave CSS Optimizer comments
		* [ ] Leave header comment
		* [ ] Leave legal comments
	* [ ] Remove unnecessary characters
		* [ ] whitespace - newlines
		* [ ] whitespace - spaces
		* [ ] whitespace - tabs
		* [ ] `;` after the last declaration within a block
	* [ ] Shorten color references to #RGB or #RRGGBB (when necessary)
	* [ ] Caching: to avoid unnecessary re-optimization
	* [ ] GZip Compression
* Configuration
	* Environment Variable to specify which options to utilize

All of these will be optional to conform with the needs of your project or task


Rational
--------

I love Go. And I've wanted to make a CSS preprocessor in it for years. A tool like [Myth][Myth (WayBack Machine)] but that didn't require an interpreter you'd have to maintain just to run the code. Just download the single binary and run it forever. Also a library that anyone can be added to their own Go project.


Inspiration
-----------


### Myth

<small>CSS the way it was imagined</small>

I loved this tool. It did everything I wanted at the time to help with browser compatibility and effectively enabling CSS vars when that was hardly supported by any browsers at the time. It was the only JavaScript based CSS preprocessor that was based on the CSS spec instead of their own language. And it was pretty fast as well.

* Language: NodeJS
* [Website][Myth (WayBack Machine)] (WayBack Machine)


### CSS-Crush

<small>CSS preprocessor</small>

When I was primary using PHP on the back-end I loved the idea of a PHP based CSS preprocessor

* Language: PHP
* [Website][CSS-Crush]



[Myth (WayBack Machine)]: https://web.archive.org/web/20201016205257/http://www.myth.io/
[CSS-Crush]: https://the-echoplex.net/csscrush/
[sindresorhus/gulp-myth: \[DEPRECATED\] Myth - Postprocessor that polyfills CSS]: https://github.com/sindresorhus/gulp-myth
[segmentio/myth: A CSS preprocessor that acts like a polyfill for future versions of the spec.]: https://github.com/segmentio/myth/
[Is this project still maintained? · Issue #150 · segmentio/myth]: https://github.com/segmentio/myth/issues/150


