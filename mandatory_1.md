## Legacy codebase dependency graphs

### Before upgrading Python and dependencies
![fb59e541-6bbe-4be4-8b1f-47638736bc83](https://github.com/user-attachments/assets/83eb04a8-2354-466c-932e-47d9abba3838)

### After upgrading Pythong and dependencies
![f8be1849-b6e5-49c3-820e-a42d3678abf7](https://github.com/user-attachments/assets/21f32f2c-7587-49f9-9954-c39e331d753a)


## Problems with legacy codebase (listed by severity high->low)
1. Outdated Language and dependencies 

2. admin password in clear text in codebase 

3. SQL injection 

4. Hashing with MD5 

5. Hardcoded configuration values 

6. no CSFR protection 

7. Single file 


## Branching strategy 
- What version control strategy did you choose and how did you actually do it / enforce it? 

  - Trunk based GitHub flow. 

  - When starting to work on a new issue, we create a branch for that issue. When done we create a pull request to main from the issue branch. All branches are based on trunk(main) and merged through pull requests directly to main on completion. When done reviewing and merging a pull request we delete the branch. 

- Why did your group choose the one you did? Why did you not choose others? 

  - We made a pros/cons for each branching strategy and figured that trunk based GitHub flow was the right one for us. As it is relatively low complexity (fewer active branches) and allows for rapid integration/regular deployment which results in getting feedback more often/faster. 

  - Some of the disadvantages to trunk based GitHub flow is that it requires discipline/thorough testing/review of new changes before they are merged into trunk. This is achieved through GitHub flow by using pull requests that require a separate reviewer to accept the changes, in addition to workflows that run tests + static code analysis on pull requests. 

  - The reason we did not choose others, such as feature/Release branching/Git flow, was due to the higher complexity introduced by having multiple long-lived branches, which can provide some good isolation/control for larger teams and projects, but we felt it was overkill for us.  

- What advantages and disadvantages did you run into during the course? 

  - Advantages: close to no merge conflicts. We have a very nice overview over what is currently being worked on and as soon as an issue is done the others can benefit from it. Enforcing only using peer reviewed PRs forces us to dive into what the others have made. 

  - Disadvantages: Using only peer reviewed PRs to merge to trunk, makes making small changes/debugging workflows tedious. 
  Feel free to revise the document even after deadline with new insights during the course. 

 
