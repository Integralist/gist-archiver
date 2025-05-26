# Project Planning #project #architecture #design #planning

## Project Planning.md

There are three types of documents you should write, in this specific order:

1. **Project** document
2. **Discovery** document
3. **Design** document

A "project" document is for planning the overarching/broader project work.\
It's responsible for breaking down the project into smaller milestones.

A "discovery" document is for _discovering_ how we'll solve ONE of the milestones defined in the project document.

A "design" document describes our design plan for what we're going to build from the discovery document.\
It will be presented to lead engineers (and architects) to review and approve.

> [!TIP]
> Also read:\
> https://www.integralist.co.uk/posts/project-management/

## Project Document

- **Preliminary Milestones:** This describes significant ‘stages’ in the project's timeline and serves as markers of progress.
- **Key Objectives:** This describes the high-level outcomes that are desired.
- **Current Architecture:** Visualise the architecture at a high-level.
- **Proposed Architecture:** Visualise the new architecture at a high-level.
- **Design Benefits:** 
- **Functional Requirements:** This defines the specific _behaviors_, _actions_, or _functions_ that the software must perform.
- **Non-functional Requirements:** This describes the _characteristics_ of the system (e.g. performance, security, reliability etc).
- **Miscellaneous:** This section includes miscellaneous documents and links relevant to the overall discussion.

## Discovery Document

This document presents an objective assessment of the existing system's requirements and deficiencies, and proposes a high-level approach that will inform the creation of a more detailed design document.

- **Glossary:** Make sure everyone is on the same page as to what certain words mean.
- **Customer Abstraction:** What the customer perceives is happening. It’s what they buy/use/derive value from.
- **Description:** What is the current situation and what are the concerns?
- **Problem:** What are the underlying problems contributing to this situation?
- **Approach:** What will be our approach to resolve these issues?
- **Business Outcomes:** The desired, measurable results of this project that directly benefit the business.

## Design Document

- **Goals:** What specifically are we trying to accomplish? By when? Are there guiding principles to follow?
- **Requirements:** Include any technical (storage, accessibility, latency/speed, monitoring & alerting), security, compliance or Product-requested deliverables.
- **Out of Scope (Non-Goals):** What are we thoughtfully and purposefully excluding from this project?
- **Success metrics:** How will we know when the project is done? Include any timelines to meet.
- **Proposed Design:** How are we proposing to solve the problems: what are we doing, what systems will change?
- **Rationale:** Why did we choose this approach? This should include a link to [a matrix spreadsheet](https://docs.google.com/spreadsheets/d/1ZnxIY4BCnsUaY65Cc2GCmzmq18XBJvG2iErSWfMjW10/edit?usp=sharing).
- **Interactions with existing systems:** What systems/components are available to build from or interact with? Are there relationships to common Fastly architecture principles?
- **Resources:** What resources do we need for this proposed design? Are there licenses, versions, hardware, specific expertise and training,etc to include?
- **Development Stages:** Identify the implementation stages of the project, so people understand what roll out looks like and how it will be approached.
- **Alternatives considered:** Describe other options considered and trade off/decisions for not doing.
- **Considerations, Risks or Constraints:** Add any special considerations, risks, or constraints. 
- **Open Questions:** Note anything for further investigation or future problems to solve (create action items and assign for follow-up if needed).
- **Stakeholders & Approval:** Primary stakeholders listed check the box to indicate sign-off or add comments with questions/concerns.
- **References:** Include any links to related documents: Problem Statement, PRD, POC proposals, boards or architecture documentation for additional context.

