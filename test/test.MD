# 89minutes
A simple multi-user streaming and blogging platform, allowing users to stream videos, publish and manage the blogs under their personal publication and readers to read them.

## Core Requirements

 - The 89minutes to be able to go live with topic of the day and stream in 4K+ resolution.
 - The users should be able to go live and share their thoughts.
 - Writers should be able to **publish** blog under their personal publication.
 - Readers should be able to **read** the blog, like and comment on the posts, bookmark, contribute to the author(tokens and vCAsh)  
 - A user can be both - a reader as well as a writer
 - Author of the blog should be able to **delete** the blog if the ownership isn't given to 89minutes.
 - Blog may **contain images and videos**, but will not contain other file types.
 - Time to access the blog should be **as low as possible**
 - We have to render "**number of blogs and streamed videos**" written by every user on his/her profile.
 - Users should be able to **search** for a particular blog
 - The platform should be scaled for **5 million** daily active readers
 - The platform should be scaled for **10,000** daily active writers


## Topics
1. Paragraphs 

    MD expects a full line space to show texts in a different line else it joins text in the same line.
2.  Text decorations

    MD can write **bold** texts, ~~italiic~~ *italic*  texts
3. Headings
    No of #'s represent the type of heading. Github will automatically add id's to headings, so the text will be automatically linked. 
    ## This is h2
    ### This is h3
4. Links

   [My Github](https://github.com/bhupendra1011 "all repos") account.[Bhupendra][1] github repo.

5. Images
    Images can be used just like links. ![Alt txt](img url)

    !["cat Img"](img2.jpeg)

    Thumbnails images can also be used which links to larger image 
    [<img src="img1.jpeg">](http://placekitten.com/200/200)

6. Ordered and Unordered Lists

    Coding Best Practices:

    * Keep code DRY
    * Writing Unit Test cases
    * Checking cross-browser support

    Steps to merge branch:

    1. Create a branch from feature
    1. commit your changes
    1. push your changes
    1. raise a pull request

7. Code Blocks

    This is super helpful when posting any code snippet


    ```js
    const fn = () => alert("some fn");
    ```




    ```css
    .hide {
        display:none
    }
    ```


    Also can show code difference


    ```diff
    var x = 10;
    - const counter = 0;
    + let counter = 0
    ```



8. Tables 

    Tables can be generated with headings and text alignment option

    |Stocks|Price|
    |:-----:|------:|
    |TCS|230|
    |YES Bank|500|



Cool Tips 

 * [Grammarly](https://marketplace.visualstudio.com/items?itemName=znck.grammarly) extension can eliminate typo and grammar mistakes
 * [ScreenTOGif](https://www.screentogif.com/) to record videos in GIF format
 * Upload GIF's to [giphy](https://giphy.com/) to embed them into blog posts.
 * [Stackedit](https://stackedit.io/) for Markdown Editing in Browser.


