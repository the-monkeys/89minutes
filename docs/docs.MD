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
 - The platform should be scaled for **10,000** daily active writers.
 - The Users should have privilege to draft their stories.



##  High Level Requirements
<!--hs-->
 - Make your high-level components operate with **high availability**.
 - Ensure that the data in your system is **durable**, not matter what happens.
 - Define how your system would behave while **scaling-up** and **scaling-down**.
 - Make your system **cost-effective** and provide a justification for the same.
 - Describe how **capacity planning** helped you made a good design decision.
 - Think about how other services will interact with your service.
<!--he-->

##  Micro Requirements
<!--ms-->
 - Ensure the data in your system is **never** going in an inconsistent state.
 - Ensure your system is **free of deadlocks** (if applicable).
 - Ensure that the throughput of your system is not affected by **locking**, if it does, state how it would affect.
<!--me-->
  