# wordpress-xml-to-hugo

## Introduction

This is an attempt at an opinionated migration from WordPress to
[Hugo](https://hugo.io/). The migration is based on parsing a blog's
XML export using a parser based on 
[WordPress XML Parser](https://github.com/grokify/wordpress-xml-go).

The original parser ignores comments, therefore first
[a fork wss used](https://github.com/amanessinger/wordpress-xml-go),
and finally the parser was incorporated into this project. 

The goal is, to migrate to a completely static blog, that can be 
deployed to the free tier of [Netlify](https://www.netlify.com/).

No data about visitors will be collected, no third-party tools will
be used.

The project is motivated by the intention to migrate my 
long-standing blog from WordPress to a static site generator.

If you're interested to use this tool, I'll explain how at the end 
of this text. Skip my explanations if you will, but I strongly
advise you to read through once. I'll explain my options and choices,
and they may turn out to be different from yours.

## Rationale for the choice of tools

Last time I checked (mid-June 2018), 
[StaticGen](https://www.staticgen.com/) listed 224 options for 
static generators. Maybe a third is not intended to be used for blogs,
but this still leaves more than I was willing to try. 

* The initial idea was, to select a tool that's written in one of the
  languages I have been working with recently. This would be
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
React applications, Jekyll and Hexo have a migration tool. As I
was neither proficient in Ruby nor Go, Hexo was a clear choice
for first experiments.

It didn't end well. Hexo had no problem migrating 4300 posts to
Markdown, but upon rendering, it didn't even manage to load all those
posts. Without generating a single HTML file, it crashed after an 
hour with out-of-memory.

To be fair, the test was executed in a RHEL 7 VM with only one
CPU and 1 GB of memory. Hexo might very likely have fared better 
when running directly on my MacBook Pro, but nevertheless,
the bad impression was there. 

Hugo was supposed to be the speed king when rendering. In order to 
get a feeling of what can possibly be expected, I installed
Hugo (one binary only) directly on the Mac. Hugo rendered the whole
blog in maybe 10 seconds and the result was already a blog, linking,
categories and tags included. It also looked gorgeous. What should
have been an experiment, solely for orientation, turned out to be 
a mighty incentive for reconsidering the selection criteria.

This time the approach was backward. [Hugo](https://hugo.io/) had to
be the winner. The idea now was, to find out if there was a viable
migration path.

* As to the language [Go](https://golang.org/), there had already been
  interest in learning it. I had only lacked a good reason to
  do so. Therefore not knowing the language turned out to be
  a welcome excuse to learn it.

* Hugo has a converter based on a WordPress plugin, but it did not
  produce the expected results of my blog. The Hugo docs
  recommend the Jekyll converter as a fallback, but although the
  preservation of HTML produced results pleasingly similar to the 
  original layouts, the blog structure was not properly preserved.
  Actually now I know that this was my fault. Hugo themes are either
  written to expect posts under "content/post" or under 
  "content/posts". If you apply the wrong theme to a your chosen
  structure, no posts are found. Anyway. There were other gripes
  as well, for instance the choice of file name and project 
  structure. 4300 posts in one directory, that's not very usable.

* The Jekyll converter did extract comments (which Hexo had not), 
  but they were put into the Markdown file's front matter and were 
  ignored by Hugo

In any case the experiments indicated that some kind of custom 
conversion would be necessary. 

The markdown files out of Jekyll's converter did contain the necessary 
information, but they would need extensive post-processing. 
Instead of going that route, I decided to completely do it 
myself by parsing an XML export out of WordPress.

In order to learn [Go](https://golang.org/), I chose it as 
implementation language. 
[WordPress XML Parser](https://github.com/grokify/wordpress-xml-go)
was chosen as a good starting point.

## Requirements

* The main purpose is to port my photo blog 
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
  preservation of HTML over clean, uncluttered Markdown. I may 
  write proper Markdown in the future though.

* In August 2018 the blog had 5852 approved comments. Preserving 
  these comments and their thread structure is important.

* Commenting must be possible and replying to existing comments as 
  well. 

* In the interest of rendering speed, no helper applications will be 
  used. Rendering will be done entirely by Hugo. This alone rules
  out 
  [reStructuredText](https://en.wikipedia.org/wiki/ReStructuredText),
  a format that would have allowed me to get rid of HTML tags and
  still be able to center or align images. The problem is, that 
  although Hugo supports reStructuredText, it does so by firing
  up a Python script **for each post**! Multiply that by 4300 and
  you know why I dislike the idea. Should Hugo ever get proper
  in-process support of reStructuredText, I may start using it. On 
  the other hand, reStructuredText is just another language that I
  would have to learn, and it is a language that I use for nothing
  else. If you look at it this way, using Markdown with embedded
  HTML is probably the best choice for me.

## Handling of comments

This is where it gets opinionated :)

### Options

The most common strategy for providing comments on static sites seems
to be, to employ a third-party service like 
[Disqus](https://disqus.com/). This is problematic for 
several reasons:

* Those services slow down an otherwise blindingly fast site. To be
  clear, the pages can be structured to display static content
  immediately, but there is still the feeling of slowness when one
  watches the comments loading.

* Usage of those sites is adversarial to the spirit of the 
  [GDPR](https://en.wikipedia.org/wiki/General_Data_Protection_Regulation).

* Ad-free comments are a feature of paid services only. Disqus does
  offer an ad-free service for non-commercial sites running no ads 
  themselves, but unfortunately 
  [manessinger.com](https://manessinger.com/) ran Amazon ads for
  years. This was done as a way to link to music, and although the
  ads never generated income worth talking about, the code is 
  wrought into the pages. I chose to filter the ads out during 
  conversion, so theoretically I could use Disqus for free. I
  still don't like it though, because of the last point:

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
out for my purposes. Still, if an alternative solution 
can be found, that would be strongly preferred.

One such alternative could be based upon the algorithms implemented
in [Antispam Bee](https://github.com/pluginkollektiv/antispam-bee).

Execution of antispam code would already be in a serverless function.
The outcome would be either rejection or, hopefully, submission
to me via email. In low-volume situations like mine, 
this might already suffice. I would merge the new
comment into the appropriate file and that would be it.

For higher-volume situations like that of for instance Mike
Johnston's of 
[The Online Photographer](http://www.theonlinephotographer.com/), 
merge support would be a nice-to-have feature. Mike may get 50 or
100 comments per day. In his case it does not matter, he's curating
comments anyway, but manual merging would quickly become a chore,
especially when threads should be preserved.

### Implementation

In the end I have decided to generate a separate folder "comments"
in the Hugo project's top-level directory. Its directory structure
mirrors that of "content". For every post with comments, there is
a directory containing one JSON file per comment. The comments
are not merged into the markdown files, but instead they are
read by Hugo and rendered by the theme's "single.html" template.

The approach is based on the excellent article 
[How To Have Comments In Hugo Without An External Service](http://saimiri.io/2016/06/comments-in-hugo/) by 
[Juha Auvinen](https://github.com/saimiri). His tool 
[hugojoco](https://github.com/saimiri/hugojoco) is meant to listen
on a server, therefore I didn't use it directly, but the base
ideas are his.

For any further details please look directly at the 
[repo of my blog](https://github.com/amanessinger/manessinger.com).

## How to use this tool

Congratulations, you've made it down here and you seem to be
still interested. Let's talk about using this tool.

1. Install [Go](https://golang.org/) and learn it. At least 
   play through the [Tour](https://tour.golang.org/). You
   should also try to understand 
   [Go Templates](https://golang.org/pkg/text/template/) as I use 
   them in the converter, and so does Hugo for templating. You
   need that knowledge, or otherwise you'll have a hard time
   adapting Hugo themes to your needs.

2. Fork this project. You will need to make changes for your 
   exact situation. No need to create pull requests, I won't
   incorporate those changes. If you find bugs though, a pull
   request would be perfectly appropriate.

3. [pkg/converter/config.go](pkg/converter/config.go) is the 
   one file that you **must** adapt.

4. Test your changes, change to the directory 
   [wordpress-xml-to-hugo](wordpress-xml-to-hugo) and 
   `go install` the program.

5. Have an XML export of your WordPress blog ready.

6. Install [Hugo](https://gohugo.io) and use it to create a new
   project for your blog. Use `git init` to add it to version
   control
   
7. Install a [Hugo theme](https://themes.gohugo.io/). I've used 
   [Natrium](https://themes.gohugo.io/hugo-natrium-theme/), 
   checked it out in my blog project's "themes" folder 
   (not as a Git sub-project), removed the `.git` subdirectory of
   the theme and instead added the theme to my blog project's repo.
   This way I lose any future theme improvements, but I need to
   make big changes anyway, for instance for rendering comments.
   
8. Make sure that your theme expects posts under the path defined
   as `PostDirectoryContentSubPath` in 
   [pkg/converter/config.go](pkg/converter/config.go). Otherwise 
   change that definition.

9. Go into that project's top-level directory and from there run 
   
       wordpress-xml-to-hugo <path-to-wp-export> .

10. Use `hugo server` to inspect the result. You might need to
   look at a lot of posts until you can be sure that the 
   conversion is acceptable. Some iterations might be needed.

11. When you're sure that the result is as good as an automatic
   conversion can reasonably be, commit the converted result to 
   version control. All further changes will be manual.

12. Decide about a way to display existing comments and to handle
   future ones. Look at
   [my blog](https://github.com/amanessinger/manessinger.com) for
   an example. As of 2018-08-15, display of comments is done,
   enabling new comments is work in progress.
   
## Conclusion

As far as I'm concerned, this project is done. I have achieved my 
goals, the blog has successfully been converted, all information
that I care about has been ported over.