# Mandatory II

## Reflection on usage of version control
Hver gang vi starter nyt issue op, laver vi en ny branch fra main med følgende flow:
1. `git switch main` / `git checkout main` - hvis man ikke allerede er i main.
2. `git pull` - for at sikre sig at main er opdateret.
3. `git switch -c <issue_id-description>` / `git checkout -b <issue_id-description>`. - opret og skift til den nye branch.

Herefter arbejder vi i branchen hvor vi løbende comitter og pusher kode til branchen på remote.
1. `git push -u origin <branch_name>` - pusher branchen til remote.
2. For at få et bedre overblik over hvilke filer vi stager benytter vi vscodes indbyggede source control.
3. `git commit -m "<issue_id description>` - commit med issue id og besked
4. `git push` - pusher commits til remote

Når vi er færdige med issuet opretter vi en pull request via GitHubs pull request extension til vscode.

## How are you DevOps?

### The Principles of Flow

#### ***Make Work Visible***
Alle opgaver er organiseret som issues på vores GitHub Projects Kanban-tavle, hvilket sikrer gennemsigtighed for alle teammedlemmer. Automatiserede release-noter opsummerer nye funktioner og opdateringer, mens obligatoriske peer-reviews for pull requests holder hele teamet informeret om ændringer og involveret i udviklingsprocessen.

#### ***Limit Work in Progress***
Hvert teammedlem arbejder typisk kun på én opgave ad gangen. For at opretholde en effektiv fremgangsmåde opdeler vi funktioner i de mindst mulige opgaver (issues), hvilket hjælper os med at fokusere og undgå multitasking.

#### ***Reduce Batch Sizes***
Ved at opdele funktioner i mindre opgaver, søger vi at holde hver arbejdsmængde håndterbar, så det bliver lettere at afslutte, gennemgå og integrere.

#### ***Reduce the Number of Handoffs***
Selvom vores obligatoriske peer-review-proces naturligt introducerer handoffs, minimeres disse ved at sikre, at hver pull request (PR) er lille og let at gennemgå. Denne tilgang reducerer kompleksiteten af handoffs, selvom vi har hyppige, lette PR'er.

#### ***Continually Identify and Evaluate Constraints***
Vi vurderer løbende vores processer for flaskehalse eller tidskrævende manuelle opgaver. Når vi støder på disse, arbejder vi på at omarrangere eller optimere flowet for at reducere ventetider og forbedre effektiviteten, hvor det er muligt.

#### ***Eliminate Hardships and Waste in the Value Stream***
Vi stræber efter at automatisere alle gentagne trin med henblik på at reducere tidsspild og maksimere effektiviteten. Ved at automatisere workflows og optimere rækkefølgen af trin i vores CI/CD pipeline udnytter vi vores ressourcer bedst muligt og undgår unødvendige forsinkelser.

### The Principles of Feedback

#### ***See Problems as They Occur***
Automatiserede tests på pull request-stadiet gør det muligt for os at fange fejl tidligt, inden de når main-branchen. Derudover kører Postman-monitoren flere gange dagligt på vores produktion, og sender os e-mails, hvis der opstår problemer, så vi hurtigt kan tage hånd om dem.

#### ***Swarm and Solve Problems to Build New Knowledge***
Når der opstår hændelser, som f.eks. nedetid, informerer vi straks gruppen. Tilgængelige teammedlemmer samarbejder for at identificere grundårsagen og udvikle en løsning, hvilket styrker vores fælles læring og vidensdeling.

#### ***Keep Pushing Quality Closer to the Source***
Vi tilstræber at løse problemer og sikre kvaliteten så tidligt som muligt, enten gennem grundige peer-reviews eller proaktiv overvågning i pipeline’en.

#### ***Enable Optimizing for Downstream Work Centers***
Dette princip er ikke aktuelt for os, da vi ikke har downstream arbejdscentre.

### The Principles of Continual Learning and Experimentation

#### ***Institutionalize the Improvement of Daily Work***
Selvom vi ikke arbejder på projektet dagligt, søger vi løbende måder at forbedre vores samarbejde og workflows på, for at sikre en konstant fremgang i teamets dynamik og produktivitet.

#### ***Transform Local Discoveries into Global Improvements***
Når en af os finder en innovativ løsning eller effektivisering, deler vi den med teamet, så alle kan få gavn af disse indsigter og vores fælles viden vokser.

#### ***Inject Resilience Patterns into Our Daily Work***
Vores engagement i at lære af hændelser og forbedre vores processer hjælper os med at indbygge robusthed i vores workflow, så vi kan håndtere udfordringer mere effektivt over tid.

#### ***Leaders Reinforce a Learning Culture***
Uden formelle ledere arbejder vi med en samarbejdende tilgang, hvor alle medlemmer opfordres til at lære og udvikle sig. Vores tilgang anerkender, at ingen forventes at vide alt fra starten; vi er alle en del af en løbende læringsproces.

## Analyze software quality
Vi implementerede SonarQube scan som en del af vores pipeline på både pull requests og på merges til main branchen tidligt i forløbet. Sidenhen har vi tilføjet OWASP ZAP på såvel frontend som backend ved pull request. Desuden har vi tilføjet GitHubs CodeQL analysis der dels afvikles på events såsom pull requests og merge, men også som cron job hvor hele repositoriet scannes.

***SonarCloud:***

SonarCloud og SonarLint extension til vscode har over flere omgange hjulpet os med at reducerer duplikeret kode, flagge usikker brug af regEx, reducerer kognitiv komplexitet og blive klogere på best practices.
Et konkret eksempel herpå er brug af read-only props i React

Det har hjulpet os til at blive vidende om forskellige best practices vi ikke var bekendte med. Det har tvunget os til at være mere generiske end vi nok ellers ville have været, hvilket har en positiv indflydelse på vedligeholdelse af kodebasen.

***OWASP ZAP:***

Ved at benytte OWASP ZAP er vi blevet bekendt med og har forsøgt at løse forskellige problematikker som blev flagget med medium severity. Et konkret eksempel er web cache desception hvor vi først og fremmest måtte undersøge hvad termet dækker over og hvordan det fungerer, og herefter finde ud af hvad vi konkret kunne implementerer af foranstaltninger for at minimere risikoen.

***GitHub CodeQL:***

GitHubs CodeQL gav os et hav af fejl og foreslag til rettelser i forbindelse med implementering af logging på backenden, dels i forhold til at logge ting såsom brugernavne, emails, tokens med mere, og dels i forhold til sanitering af brugerinput, som ikke var et koncept vi hidtil var blevet introduceret til på uddannelsen.

CodeQLs måde at flagge issues direkte i pull requestet som kommentarer og at den giver forslag til ændringer fra co-pilot som kan committes direkte er en effektiv feature, men den kan være uoverskuelig når man ikke helt forstår problemet den forsøger at løse som eksempelvis sanitering af brugerinput.

## Monitoring realization

I forbindelse med at vores Azure VM crashede relativt tidligt i forløbet benyttede vi Azures indbyggede monitorering til at finde årsagen til crashet, hvilket var at vores VM løb tør for memory.
Derfor opskalerede vi til det billigste betalte niveau og har ikke siden haft problemer.
Dette har måske bidraget til at vi ikke som sådan har fundet noget alarmerende i forbindelse med vores egen monitorering, men det tvinger os ligeledes til at genoverveje hvorvidt de metrics vi måler på er relavante.