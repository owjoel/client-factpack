from autogen import Agent, UserProxyAgent, GroupChat, GroupChatManager

class MainAgent(Agent): # Main Agent (Orchestrator)
    def __init__(self):
        super().__init__("MainAgent")

    def receive(self, message):
        if message["from"] == "ProcessingAgent":
            self.send_message("MLSimilarityAgent", message["data"])
            self.send_message("ResearchAgent", message["data"])
        
        elif message["from"] in ["MLSimilarityAgent", "ResearchAgent"]:
            # Store results from ML & Research Agents
            self.store_data(message["from"], message["data"])
            
            if self.has_received_from("MLSimilarityAgent") and self.has_received_from("ResearchAgent"):
                # Compare results
                ml_result = self.get_stored_data("MLSimilarityAgent")
                research_result = self.get_stored_data("ResearchAgent")

                if ml_result["match"] == research_result["match"]:
                    self.send_message("User", {"decision": ml_result, "status": "Match Confirmed"})
                else:
                    self.send_message("User", {"decision": self.make_final_decision(ml_result, research_result), "status": "MainAgent Override"})

    def make_final_decision(self, ml_result, research_result):
        # Implement logic to choose the best result
        if ml_result["confidence"] > research_result["confidence"]:
            return ml_result
        return research_result


class ProcessingAgent(Agent): # Processing Agent
    def __init__(self):
        super().__init__("ProcessingAgent")

    def receive(self, message):
        structured_data = self.process_article(message["data"])
        self.send_message("MainAgent", structured_data)

    def process_article(self, article):
        # Extract key identity fields (name, DOB, address, etc.)
        return {"name": "John Doe", "dob": "1990-01-01", "address": "123 Main St"}


class MLSimilarityAgent(Agent): # ML-Based Similarity Agent
    def __init__(self):
        super().__init__("MLSimilarityAgent")

    def receive(self, message):
        identity_match, confidence = self.check_database(message["data"])
        self.send_message("MainAgent", {"match": identity_match, "confidence": confidence})

    def check_database(self, structured_data):
        # Dummy check against identity database
        return "John Doe", 0.95  # Example: Found a match with 95% confidence


class ResearchAgent(Agent): # Research Agent
    def __init__(self):
        super().__init__("ResearchAgent")

    def receive(self, message):
        research_match, confidence = self.research_identity(message["data"])
        self.send_message("MainAgent", {"match": research_match, "confidence": confidence})

    def research_identity(self, structured_data):
        # Simulated research process
        return "John Doe", 0.85  # Example: Found match with 85% confidence


user_agent = UserProxyAgent(name="User") # User Agent (Final Decision)

# Initialize Agents
main_agent = MainAgent()
processing_agent = ProcessingAgent()
ml_agent = MLSimilarityAgent()
research_agent = ResearchAgent()

# Group Chat Setup
group_chat = GroupChat(
    agents=[main_agent, processing_agent, ml_agent, research_agent, user_agent]
)

chat_manager = GroupChatManager(group_chat=group_chat)

# Simulate a new identity matching request
main_agent.send_message("ProcessingAgent", {"data": "Article about John Doe"})
