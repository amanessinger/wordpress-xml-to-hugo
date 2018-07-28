# wordpress-xml-to-hugo

## Introduction

This is an attempt at an opinionated migration from WordPress to
[Hugo](https://hugo.io/). The migration is based on parsing a blog's
XML export using a parser based on 
[WordPress XML Parser](https://github.com/grokify/wordpress-xml-go).

The original parser ignores comments, therefore 
[a fork is used](https://github.com/amanessinger/wordpress-xml-go). 

The goal is, to migrate to a completely static blog, that can be 
deployed to the free tier of [Netlify](https://www.netlify.com/).

No data about visitors will be collected, no third-party tools will
be used.

The project is motivated by the intention to migrate the author's 
long-standing blog from WordPress to a static site generator.

## Rationale for the choice of tools

Last time the author checked (mid-June 2018), 
[StaticGen](https://www.staticgen.com/) listed 224 options for 
static generators. Maybe a third is not intended to be used for blogs,
but this still leaves more than the author was willing to try. 

* The initial idea was, to select a tool that's written in one of the
  languages the author has been working with recently. This would be
  JavaScript, TypeScript, Python or Java. Exotic languages like Haskell
  or Erlang were ruled out, legacy languages like Perl were deemed too
  old to base a new system on.

* The next criterion was popularity. A tool should have a considerable
  active user base.

* The final initial selection criterion was the existence of a 
  migration tool.

In June 2018 the top five on [StaticGen](https://www.staticgen.com/)
were Jekyll (Ruby), Next (JavaScript), Hugo (Go), Gatsby (JavaScript)
and Hexo (JavaScript). Next and Gatsby are more intended for building
React applications, Jekyll and Hexo have a migration tool. As the
author was neither proficient in Ruby nor Go, Hexo was a clear choice
for first experiments.

It didn't end well. Hexo had no problem migrating 4300 posts to
Markdown, but upon rendering, it didn't even manage to load all those
posts. Without generating a single HTML file, it crashed after an 
hour with out-of-memory.

To be fair, the test was executed in a RHEL 7 VM with only one
CPU and 1 GB of memory. Hexo might very likely have fared better 
when running directly on the author's MacBook Pro, but nevertheless,
the bad impression was there. 

Hugo was supposed to be the speed king when rendering. In order to 
get a feeling of what can possibly be expected, the author installed
Hugo (one binary only) directly on the Mac. Hugo rendered the whole
blog in maybe 20 seconds and the result was already a blog, linking,
categories and tags included. It also looked gorgeous. What should
have been an experiment, solely for orientation, turned out to be 
a mighty incentive for reconsidering the selection criteria.

This time the approach was backward. [Hugo](https://hugo.io/) had to
be the winner. The idea now was, to find out if there was a viable
migration path.

* As to the language [Go](https://golang.org/), there had already been
  interest in learning it. The author had only lacked a good reason to
  do so. Therefore not knowing the language turned out to be
  a welcome excuse to learn it.

* Hugo has a converter based on a WordPress plugin, but it did not
  produce the expected results of the author's blog. The Hugo docs
  recommend the Jekyll converter as a fallback, but although the
  preservation of HTML produced results pleasingly similar to the 
  original layouts, the blog structure was not properly preserved.

* The Jekyll converter did extract comments (which Hexo had not), 
  but they were put into the Markdown file's front matter and were 
  ignored by Hugo

In any case the experiments indicated that some kind of custom 
conversion would be necessary. 

The markdown files out of Jekyll's converter did contain the necessary 
information, but they would need extensive post-processing. 
Instead of going that route, the author decided to completely do it 
himself by parsing an XML export out of WordPress.

In order to learn [Go](https://golang.org/), it was chosen as 
implementation language. 
[WordPress XML Parser](https://github.com/grokify/wordpress-xml-go)
was chosen as a good starting point.

## Requirements

* The main purpose is to port the author's photo blog 
  [manessinger.com](https://manessinger.com/) over to a git-based 
  static workflow. This is, what the tool will be optimized for. If 
  it can be generalized for broader applicability, that's fine. No
  compromises will be made in this regard though.
  
* The blog [manessinger.com](https://manessinger.com/) had 4300 posts 
  in August 2018. The blog has had about one post per day since 
  fall 2006. Posts with more than one image always had one big,
  centered image right after the title, followed by text intermixed 
  with thumbnails of the other images. A good example is
  [here](https://manessinger.com/2018/05/4214-the-congress-hall-from-afar-v.html). This layout
  can't reasonably be changed, thus the migration needs to favor 
  preservation of HTML over clean, uncluttered Markdown.

* In August 2018 the blog had 5852 approved comments. Preserving 
  these comments and their thread structure is important.

* Commenting must be possible and replying to existing comments as 
  well. 

* In the interest of rendering speed, no helper applications will be 
  used. Rendering will be done entirely by Hugo.

## Handling of comments

This is where it gets opinionated :)

The most common strategy for providing comments on static sites seems
to be, to employ a third-party service like 
[Disqus](https://disqus.com/). This is problematic for 
several reasons:

* Those services slow down an otherwise blindingly fast site. To be
  clear, the pages can be structured to display static content
  immediately, but there is still the feeling of slowness when one
  watches the comments to load.

* Usage of those sites is adversarial to the spirit of the 
  [GDPR](https://en.wikipedia.org/wiki/General_Data_Protection_Regulation).

* Ad-free comments are a feature of paid services only. Disqus does
  offer an ad-free service for non-commercial sites running no ads 
  themselves, but unfortunately 
  [manessinger.com](https://manessinger.com/) ran Amazon ads for
  years. This was done as a way to link to music, and although the
  ads never generated income worth talking about, the code is 
  wrought into the pages. It might be possible to filter the ads out
  upon migration though.

* Migration of existing comments to such a system may likely not be
  possible, because they all rely on the users having accounts with 
  them. 

For all these reasons the comments will be rendered into the pages.
In case of existing comments this seems to be the only viable option,
and for new comments it will need a submit/approve/merge mechanism.

Submission will be through a 
[Web Component](https://en.wikipedia.org/wiki/Web_Components) 
created with [Stencil](https://stenciljs.com/). The component will
only be loaded upon click on a comment "button" rendered via SVG.
This should eliminate all SPAM submitted by bots not running in
a browser.

The Web Component would present a form and it would submit the 
comment to a serverless function. This could be anything, and the
obvious option is AWS Lambda, because the free tier of 
[Netlify](https://www.netlify.com/) already includes 125.000 free
invocations per month. This should easily cover real comments and
all SPAM possibly entered by humans.

As to SPAM filtering, WordPress uses (and provides for third 
parties to use) the Akismet service. In years of usage, Akismet has
demonstrated excellent selectivity and high reliability. It is only
not easily compatible with the 
[GDPR](https://en.wikipedia.org/wiki/General_Data_Protection_Regulation)
There are indeed doubts as to its legality at all.

Technically Akismet could easily be integrated into any comment 
submission mechanism, so its usage has not entirely been ruled
out for the author's purposes. Still, if an alternative solution 
can be found, that would be strongly preferred.

One such alternative could be based upon the algorithms implemented
in [Antispam Bee](https://github.com/pluginkollektiv/antispam-bee).

Execution of antispam code would already be in a serverless function.
The outcome would be either rejection or, hopefully, submission
to the author via email. In low-volume situations like that of the
author, this would already suffice. The author would merge the new
comment into the appropriate markdown file and that would be it.

For higher-volume situations like that of for instance Mike
Johnston's of 
[The Online Photographer](http://www.theonlinephotographer.com/), 
merge support would be a nice-to-have feature. Mike may get 50 or
100 comments per day. In his case it does not matter, he's curating
comments anyway, but manual merging would quickly become a chore,
especially when threads should be preserved.

No decision about merge support has been made yet, but it is clear 
that the submission email can easily contain all necessary 
information.

Once the comment has been merged (manually or automatically) and the 
change been pushed to GitHub, Netlify's continuous delivery should 
kick in automatically, call Hugo and let it render the changes.