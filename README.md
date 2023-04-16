## MazeMapper： A Full-Capability Recursive Resolver
MazeMapper is a powerful domain resolution tool that can perform comprehensive resolution of individual domain names, generating a complete resolution dependency graph for the domain. By analyzing the resolution dependency graph, we can identify domain resolution failures or potential causes of domain resolution failures, and then make targeted repairs. In addition, we use coroutines to resolve batches of domain names, accelerating the resolution process.


## **Environment**
MazeMapper is developed using Go 1.19.4. Before running the program, you need to install Graphviz and add it to the environment variables, so that the visualized domain resolution dependency graph can be generated.

## **Statement**
We take **<domain,qtype,Ip>** as the key value of each node in the graph, different colors represent different types of nodes.
**1.** **Green：** CNAME
**2.** **Yellow：** NS not glue IP
**3.** **Blue：** Answer   A/AAAA
**4.** **Red：** Error
For different errors, we explicitly mark them in the nodes，such as **Timeout**,**Refused**,**NameError**,**Corrupt**,**IPerror**,**NotImplemented**,**IDMisMatch**,**NoNsrecord** and so on.
