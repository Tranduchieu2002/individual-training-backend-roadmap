# 🚀 FUNDAMENTAL ENGINEER TRAINING PROGRAM
## Lộ trình 12 Tuần - BigTech Standard

---

## 📋 OVERVIEW

**Mục tiêu:** Đào tạo engineer có khả năng:
- Ship production-ready code độc lập
- Design và implement features end-to-end
- Debug và optimize hệ thống phức tạp
- Collaborate hiệu quả trong teamf**Weekly Milestone:**

**Format:**
- 40 giờ/tuần (8h coding + review + meetings)
- 1-on-1 với mentor mỗi tuần
- Code review hàng ngày
- Weekly demo & retrospective

---

# PHASE 1: FOUNDATIONS (Tuần 1-4)
## 🎯 Goal: Master clean code, testing, và Git workflow

### **WEEK 1: Clean Code & Testing Fundamentals**

#### Day 1-2: Clean Code Bootcamp
**Morning Session (4h):**
- Đọc: Clean Code Chapter 1-3 (Robert C. Martin)
- Workshop: Code smell identification
- Refactoring exercises

**Assignment:**
```
TASK: Refactor Legacy Code
- Clone: github.com/your-org/legacy-payment-service
- Branch: feat/your-name/refactor-payment-validator
- Refactor: PaymentValidator class (500 lines → <100 lines per method)
- Requirements:
  ✓ Extract methods có tên self-explanatory
  ✓ Apply Single Responsibility Principle
  ✓ Remove magic numbers
  ✓ Add meaningful variable names
  
Deliverable: PR với description giải thích từng refactoring decision
Review criteria: Readability, maintainability, no functionality change
```

#### Day 3-5: Unit Testing & TDD
**Morning Session:**
- TDD demonstration: Red → Green → Refactor cycle
- Testing patterns: AAA (Arrange, Act, Assert)
- Mocking & stubbing strategies

**Assignment:**
```
TASK: Build Feature với TDD
Feature: Promo Code Validator Service

Requirements:
1. Viết test cases trước (minimum 15 test cases):
   - Valid promo code → success
   - Expired code → error
   - Already used code → error
   - Invalid format → error
   - Concurrent usage → handle race condition
   
2. Implement để pass all tests
3. Achieve 90%+ coverage

Tech Stack: 
- Language: Go/Java/Python (choose one)
- Framework: Testing framework của language đó
- Coverage tool: go cover/JaCoCo/pytest-cov

Deliverable: 
- PR với tests + implementation
- Coverage report screenshot
- README giải thích test strategy

Review criteria: 
- Test coverage ≥ 90%
- Test readability (clear naming, no magic values)
- Edge cases covered
```

**Weekly Milestone:**
- [ ] 2 PRs merged (refactoring + TDD feature)
- [ ] Code review comments addressed < 3 iterations
- [ ] Pass code quality gate (linter + coverage)

---

### **WEEK 2: Git Workflow & Collaboration**

#### Day 1-2: Git Mastery
**Morning Session:**
- Git internals: commits, trees, blobs
- Branch strategies: GitFlow vs Trunk-based
- Interactive rebase workshop
- Conflict resolution scenarios

**Assignment:**
```
TASK: Complex Git Scenario Handling

Scenario 1: Merge Conflict Resolution
- Create feature branch từ main (đã outdated 50 commits)
- Rebase lên main mới nhất
- Resolve conflicts trong 3 files khác nhau
- Maintain commit history clean (squash WIP commits)

Scenario 2: Undo Mistakes
- Accidentally commit sensitive data (.env file)
- Remove from history với git filter-branch
- Force push safely

Scenario 3: Cherry-pick Hotfix
- Hotfix merged vào main
- Cherry-pick vào release branch
- Handle conflicts nếu có

Deliverable: Screen recording + written explanation
Review: Mentor verify qua Git history
```

#### Day 3-5: Code Review Culture
**Morning Session:**
- Code review best practices (Google Engineering Practices)
- Cách viết PR description (Context, Changes, Testing)
- Cách comment constructively

**Assignment:**
```
TASK: Code Review Training

Part 1: Review 5 PRs từ teammates
- Sử dụng review checklist:
  ✓ Functionality: Logic đúng chưa?
  ✓ Tests: Coverage đủ chưa? Edge cases?
  ✓ Performance: Có query N+1? Memory leak?
  ✓ Security: SQL injection? XSS?
  ✓ Readability: Code dễ hiểu không?
  
- Mỗi PR: ít nhất 3 meaningful comments
- Tone: Professional, helpful, specific

Part 2: Respond to Reviews on Your PR
- Address mỗi comment với:
  a) Explain reasoning, HOẶC
  b) Apply suggestion + reply "Done ✅"
  
Deliverable: 
- 5 reviewed PRs với quality comments
- Your PR updated dựa trên feedback

Review criteria:
- Comment quality (specific, actionable)
- Response time < 24h
- No defensive behavior
```

**Conventional Commits Workshop:**
```
Practice Writing Commits:

❌ Bad:
"fix stuff"
"update code"
"changes"

✅ Good:
"feat(payment): add retry logic for failed transactions"
"fix(promo): prevent race condition in code redemption"
"refactor(auth): extract JWT validation into middleware"
"docs(api): update OpenAPI spec for /orders endpoint"

Format: <type>(<scope>): <subject>
Types: feat, fix, refactor, test, docs, chore, perf

Assignment: Review last 10 commits, rewrite với Conventional Commits
```

**Weekly Milestone:**
- [ ] Git proficiency test passed (merge conflicts, rebase, cherry-pick)
- [ ] 5+ quality code reviews given
- [ ] PR description template mastered
- [ ] All commits follow Conventional Commits

---

### **WEEK 3: Database & Query Optimization**

#### Day 1-2: SQL Fundamentals & Indexing
**Morning Session:**
- Database indexing deep dive (B-Tree, Hash, Composite)
- Query execution plans (EXPLAIN ANALYZE)
- N+1 query problem & solutions

**Assignment:**
```
TASK: Database Optimization Challenge

Given: E-commerce database với slow queries
Schema:
- users (1M records)
- orders (5M records)  
- order_items (20M records)
- products (100K records)

Problem Queries (chạy > 5 seconds):

1. "Get all orders with items for user X in last 30 days"
   SELECT * FROM orders o
   JOIN order_items oi ON o.id = oi.order_id
   WHERE o.user_id = ? AND o.created_at > NOW() - INTERVAL 30 DAY

2. "Top 10 best-selling products this month"
   SELECT p.name, COUNT(*) as sales
   FROM products p
   JOIN order_items oi ON p.id = oi.product_id
   JOIN orders o ON oi.order_id = o.id
   WHERE o.created_at > DATE_TRUNC('month', NOW())
   GROUP BY p.id
   ORDER BY sales DESC
   LIMIT 10

Your Tasks:
1. Run EXPLAIN ANALYZE trên mỗi query
2. Identify bottlenecks (seq scan, missing index...)
3. Add appropriate indexes
4. Rewrite queries nếu cần
5. Optimize to < 100ms

Deliverable:
- SQL migration file với indexes
- Before/After EXPLAIN ANALYZE comparison
- Document: Why each index helps

Review criteria:
- Query time reduced ≥ 95%
- Index choices justified
- No over-indexing (storage/write trade-off)
```

#### Day 3-5: Transaction Management
**Morning Session:**
- ACID properties deep dive
- Isolation levels & their trade-offs
- Deadlock scenarios & prevention

**Assignment:**
```
TASK: Fix Race Condition Bug

Scenario: Promo Code System
Bug Report:
"Users can use same promo code multiple times by rapidly clicking 
'Apply Code' button. Lost revenue: $50K last month."

Current Code (pseudocode):
```go
func ApplyPromoCode(userID, code string) error {
    promo := db.Query("SELECT * FROM promo_codes WHERE code = ?", code)
    
    if promo.Used {
        return ErrAlreadyUsed
    }
    
    if promo.ExpiresAt < time.Now() {
        return ErrExpired
    }
    
    // Simulate processing delay
    time.Sleep(100 * time.Millisecond)
    
    db.Exec("UPDATE promo_codes SET used = true WHERE code = ?", code)
    return nil
}
```
Your Tasks:
```
1. Write test reproducing race condition (concurrent requests)
2. Propose 3 solutions:
   a) Pessimistic locking (SELECT FOR UPDATE)
   b) Optimistic locking (version column)
   c) Distributed lock (Redis)
   
3. Implement ALL 3 solutions in separate branches
4. Benchmark each solution (throughput, latency)
5. Write RFC comparing trade-offs
```
Deliverable:
- 3 PRs (one per solution)
- Load test results (use k6/JMeter)
- RFC document with recommendation

Review criteria:
- All solutions prevent race condition
- Performance benchmarks accurate
- RFC reasoning sound (consider scale, complexity, ops)

**Weekly Milestone:**
- [ ] Optimize ≥ 3 slow queries to < 100ms
- [ ] Fix race condition với proper locking
- [ ] Understand transaction isolation levels
- [ ] Write migration scripts safely

---

### **WEEK 4: System Observability**

#### Day 1-3: Logging & Metrics
**Morning Session:**
- Structured logging (JSON format)
- Metrics types: Counter, Gauge, Histogram
- Log aggregation (ELK stack overview)

**Assignment:**
TASK: Add Observability to Service

Service: Order Processing System
Current State: No logs, no metrics, blind when issues happen

Requirements:

1. Structured Logging:
   - Replace println() với proper logger (zerolog/zap/logrus)
   - Log format:
     {
       "timestamp": "2024-01-15T10:30:00Z",
       "level": "error",
       "service": "order-processor",
       "trace_id": "abc-123",
       "user_id": "user-456",
       "message": "Payment validation failed",
       "error": "insufficient_funds",
       "amount": 150.00
     }
   
   - Log levels strategy:
     * DEBUG: Function entry/exit
     * INFO: Business events (order created, payment processed)
     * WARN: Recoverable errors (retry attempts)
     * ERROR: Failures requiring attention
     
2. Metrics (Prometheus format):
   - Counter: orders_total{status="success|failed"}
   - Histogram: order_processing_duration_seconds
   - Gauge: active_orders_count
   
3. Distributed Tracing:
   - Add trace_id propagation
   - Span for each major operation

Deliverable:
- PR với logging + metrics instrumented
- Grafana dashboard JSON (mock với screenshots)
- Runbook: "How to debug order processing issues using logs/metrics"

Review criteria:
- Logs contain enough context for debugging
- Metrics cover key business + technical KPIs
- Dashboard useful for on-call engineers

#### Day 4-5: Debugging Production Issues
**Simulation Exercise:**
INCIDENT SIMULATION: "Orders Processing Slowing Down"

Setup:
- You're on-call
- PagerDuty alert: "P95 latency > 5s (SLA: 1s)"
- No other context

Your Tools:
- Grafana dashboards
- Kibana logs
- Database slow query log
- APM traces (Datadog/New Relic)

Tasks:
1. Investigate root cause (you have 30 minutes)
   - Check metrics: Which service is slow?
   - Check logs: Any errors?
   - Check database: Slow queries? Lock waits?
   - Check traces: Where is time spent?
   
2. Mitigate immediately:
   - Can we add cache?
   - Can we add index?
   - Can we kill slow queries?
   - Can we scale horizontally?
   
3. Write post-mortem:
   - Timeline of events
   - Root cause analysis (5 Whys)
   - Short-term fix
   - Long-term prevention
   - Action items

Deliverable:
- Incident report document
- PR với fix (if applicable)
- Runbook update

Review criteria:
- Root cause identified correctly
- Mitigation effective
- Post-mortem thorough
- Lessons learned documented


**Weekly Milestone:**
- [ ] Service fully instrumented (logs + metrics + traces)
- [ ] Create ≥ 3 useful Grafana dashboards
- [ ] Successfully debug simulated incident
- [ ] Write runbook for common issues

---

# PHASE 2: SYSTEM DESIGN (Tuần 5-8)
## 🎯 Goal: Design scalable systems end-to-end

### **WEEK 5: Architecture Patterns**

#### Day 1-2: Microservices vs Monolith
**Morning Session:**
- When to use Monolith vs Microservices
- Service boundaries design
- Inter-service communication (REST, gRPC, Events)

**Assignment:**
```
TASK: Design Migration Plan

Current: Monolithic E-commerce App (100K LOC)
Problem: 
- Deploy takes 2 hours (full regression test)
- One bug in checkout breaks entire site
- Team of 50 engineers stepping on each other
```
Your Task: Propose Migration to Microservices

Deliverable: Architecture Design Doc (ADD)

Template:
# Architecture Design Doc: E-commerce Microservices

## 1. Current State Analysis
- Draw monolith architecture diagram
- Identify coupling points
- List pain points with evidence

## 2. Proposed Architecture
- Service boundaries (e.g., User, Product, Order, Payment, Notification)
- Communication patterns:
  * Synchronous: REST/gRPC for read/write
  * Asynchronous: Kafka/RabbitMQ for events
- Data ownership (each service owns its DB)

## 3. Migration Strategy
- Phase 1: Extract low-risk service (e.g., Notification)
- Phase 2: Extract high-value service (e.g., Payment)
- Phase 3: Strangler pattern for remaining
- Timeline: 6-12 months

## 4. Trade-offs & Risks
| Decision | Pros | Cons | Mitigation |
|----------|------|------|------------|
| gRPC vs REST | Performance | Complexity | Start with REST, migrate critical paths |
| Kafka vs RabbitMQ | Scalability | Ops overhead | Use managed service (AWS MSK) |

## 5. Success Metrics
- Deploy frequency: 1/week → 10/day
- MTTR: 2h → 15min
- Team velocity: +30%

Review criteria:
- Service boundaries logical (high cohesion, low coupling)
- Migration plan realistic
- Trade-offs clearly explained
- Risks identified with mitigation


#### Day 3-5: Caching Strategies
**Assignment:**
```
TASK: Implement Multi-Layer Caching

Service: Product Catalog API
SLA: P95 latency < 50ms, 10K QPS

Current State:
- Every request hits database
- P95 latency: 300ms
- Database CPU: 80%

Requirements:
1. Design caching strategy:
   - L1: Application cache (in-memory)
   - L2: Redis (distributed cache)
   - L3: CDN (for static content)

2. Implement:
   - Cache-aside pattern
   - Write-through for updates
   - TTL strategy (short for hot items, long for cold)
   - Cache invalidation on product updates

3. Handle edge cases:
   - Cache stampede (thundering herd)
   - Stale data during deployment
   - Cache penetration (query for non-existent items)

Code Example (Go):
```go
type ProductService struct {
    cache      cache.Cache
    redis      *redis.Client
    db         *sql.DB
}

func (s *ProductService) GetProduct(id string) (*Product, error) {
    // L1: Check app cache
    if product := s.cache.Get(id); product != nil {
        return product, nil
    }
    
    // L2: Check Redis
    if product := s.redis.Get(id); product != nil {
        s.cache.Set(id, product, 1*time.Minute)
        return product, nil
    }
    
    // L3: Query DB
    product := s.db.Query(...)
    
    // Populate caches
    s.redis.Set(id, product, 10*time.Minute)
    s.cache.Set(id, product, 1*time.Minute)
    
    return product, nil
}
```

Deliverable:
- Implementation với all 3 cache layers
- Load test comparison (before/after)
- Cache hit rate monitoring

Review criteria:
- P95 latency < 50ms achieved
- Cache hit rate > 90%
- Edge cases handled properly


**Weekly Milestone:**
- [ ] Complete architecture design doc
- [ ] Implement multi-layer caching
- [ ] Reduce latency by ≥ 80%
- [ ] Understand CAP theorem trade-offs

---

### **WEEK 6-7: CI/CD & Infrastructure**

#### Week 6: CI/CD Pipeline
**Assignment:**
```
TASK: Build Production-Grade CI/CD Pipeline

Requirements:
1. Build Pipeline (GitHub Actions/GitLab CI):
   ```yaml
   stages:
     - lint
     - test
     - build
     - deploy
   
   lint:
     - Run linter (golangci-lint/eslint)
     - Check code formatting
     - Fail if issues found
   
   test:
     - Unit tests (parallel execution)
     - Integration tests
     - Coverage report
     - Fail if coverage < 80%
   
   build:
     - Build Docker image
     - Scan for vulnerabilities (Trivy)
     - Push to registry (tag với commit SHA)
   
   deploy:
     - Deploy to staging (auto)
     - Run smoke tests
     - Deploy to production (manual approval)
     - Rollback capability
   ```

2. Infrastructure as Code:
   - Kubernetes manifests (Deployment, Service, Ingress)
   - Helm chart packaging
   - Environment configs (dev/staging/prod)

3. Deployment Strategy:
   - Blue-Green deployment
   - Health checks (liveness, readiness)
   - Resource limits (CPU, memory)

Deliverable:
- .github/workflows/ci-cd.yml
- Kubernetes manifests in /k8s folder
- Deployment runbook

Review criteria:
- Pipeline runs in < 10 minutes
- Zero-downtime deployment
- Rollback works in < 2 minutes
---
#### WEEk 7: Docker & Kubernetes
**Assignment:**
TASK: Containerize & Orchestrate Service

Part 1: Optimize Docker Image
Current Dockerfile (1.2GB):
```dockerfile
FROM ubuntu:latest
RUN apt-get update && apt-get install -y golang
COPY . /app
WORKDIR /app
RUN go build -o server
CMD ["./server"]
```

Your Optimized Dockerfile:
- Use multi-stage build
- Alpine base image
- Non-root user
- Only copy necessary files
- Target: < 50MB

Part 2: Kubernetes Deployment
```yaml
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: order-service
spec:
  replicas: 3
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0
  template:
    spec:
      containers:
      - name: order-service
        image: your-registry/order-service:v1.2.3
        resources:
          requests:
            cpu: 100m
            memory: 128Mi
          limits:
            cpu: 500m
            memory: 512Mi
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
        env:
        - name: DB_HOST
          valueFrom:
            secretKeyRef:
              name: db-credentials
              key: host
```

Tasks:
- Setup local Kubernetes (minikube/kind)
- Deploy service với proper configs
- Test rolling update (deploy new version)
- Simulate pod failure (kill pod, verify auto-restart)
- Setup HPA (Horizontal Pod Autoscaler)

Deliverable:
- Optimized Dockerfile
- Complete Kubernetes manifests
- Load test showing auto-scaling

Review criteria:
- Image size < 50MB
- Zero downtime during updates
- Auto-scaling works under load
**Weekly Milestone:**
- [ ] CI/CD pipeline running smoothly
- [ ] Docker image optimized
- [ ] Service running on Kubernetes
- [ ] Deployment automation complete
---

### **WEEK 8: Performance & Scalability**

**Assignment:**
TASK: Performance Optimization Challenge

Service: Real-time Notification System
Current Performance:
- 1K messages/second
- P95 latency: 2s
- CPU: 90%

Target:
- 10K messages/second
- P95 latency: 100ms
- CPU: < 60%

Optimization Areas:

1. Code-Level:
   - Profile with pprof/py-spy
   - Identify bottlenecks (CPU, memory, I/O)
   - Optimize hot paths
   - Use connection pooling
   - Implement batching

2. Architecture-Level:
   - Add message queue (Kafka/SQS)
   - Implement consumer groups
   - Scale horizontally
   - Add Redis for deduplication

3. Database-Level:
   - Add indexes
   - Use read replicas
   - Implement write buffering

Deliverable:
- Profiling report (before/after flame graphs)
- Load test results (use k6):
  * Throughput comparison
  * Latency distribution
  * Error rate
- Architecture diagram showing changes
- Cost analysis (infra cost before/after)

Review criteria:
- Meet all performance targets
- Solution scales linearly
- No single point of failure
- Cost-effective
```

**Capstone Project Kickoff:**
```
CAPSTONE PROJECT: Build Mini Grab

Timeline: 4 weeks (Week 9-12)
Team: Solo or pair

Requirements:
1. Rider App API:
   - Register/login
   - Request ride
   - Track driver location (real-time)
   - Rate driver

2. Driver App API:
   - Accept/reject rides
   - Update location (GPS coordinates)
   - Complete ride

3. Admin Dashboard:
   - View all active rides
   - Monitor system health
   - Generate reports

Tech Stack:
- Backend: Go/Java/Python
- Database: PostgreSQL + Redis
- Message Queue: Kafka/RabbitMQ
- Real-time: WebSocket/SSE
- Infrastructure: Docker + Kubernetes
- Monitoring: Prometheus + Grafana

Non-Functional Requirements:
- API latency: P95 < 200ms
- Support 1K concurrent rides
- 99.9% uptime
- Complete observability
- CI/CD pipeline
- < $100/month cloud cost

Week 9-10: Backend implementation
Week 11: Frontend + Integration
Week 12: Load testing + Presentation

---

# PHASE 3: PRODUCT ENGINEERING (Tuần 9-12)
## 🎯 Goal: Ship features end-to-end với product mindset

### **WEEK 9-10: Capstone Implementation**

**Daily Standups:**
- What did I do yesterday?
- What will I do today?
- Any blockers?

**Mentor 1-on-1 (2x/week):**
- Code review deep dives
- Architecture decisions review
- Unblock technical challenges

**Mid-Point Review (End of Week 9):**
```
Presentation Format:
1. Demo (10 min):
   - Show working features
   - API endpoints demo (Postman/curl)
   - Database schema walkthrough

2. Architecture (10 min):
   - System diagram
   - Technology choices justification
   - Scaling strategy

3. Q&A (10 min):
   - "How do you handle race conditions in ride matching?"
   - "What happens if Kafka goes down?"
   - "How do you ensure location privacy?"

Deliverables:
- Working MVP (core features functional)
- Architecture design doc
- API documentation (OpenAPI/Swagger)
```

---

### **WEEK 11: Product Mindset & A/B Testing**

#### Feature Launch Workshop
**Morning Session:**
- The "Why" before "How"
- Feature prioritization (RICE framework)
- Writing PRDs (Product Requirement Docs)

**Assignment:**
```
TASK: Launch Feature với Data-Driven Approach

Feature: Dynamic Pricing (Surge Pricing)
Business Goal: Increase revenue by 15% during peak hours

Your Tasks:

1. Write Product Requirement Doc:
   # PRD: Dynamic Pricing v1
   
   ## Problem Statement
   During peak hours (7-9am, 5-7pm), driver supply < demand.
   Users wait 15+ minutes. Lost revenue: $1M/month.
   
   ## Proposed Solution
   Surge pricing: Price multiplier based on supply/demand ratio.
   
   ## Success Metrics
   - Revenue: +15%
   - User churn: < 5%
   - Driver acceptance rate: > 80%
   
   ## User Stories
   - As a user, I see price multiplier before booking
   - As a driver, I earn more during surge
   
   ## Edge Cases
   - Max surge cap: 3x
   - Gradual increase (not instant spikes)
   - Notify users of price changes

2. Implement Feature Flag:
   ```go
   if featureFlag.IsEnabled("dynamic-pricing", userID) {
       price = basePrice * surgeFactor
   } else {
       price = basePrice
   }
   ```

3. Design A/B Test:
   - Control: 50% users (no surge)
   - Treatment: 50% users (with surge)
   - Duration: 2 weeks
   - Monitor metrics daily

4. Build Analytics Dashboard:
   - Revenue per user (control vs treatment)
   - Booking conversion rate
   - User churn rate
   - Driver earnings

Deliverable:
- PRD document
- Feature implementation với flag
- A/B test plan
- Mock analytics dashboard

Review criteria:
- PRD clear and actionable
- Metrics well-defined
- A/B test statistically sound
- Edge cases handled
---

### **WEEK 12: Final Demo & Retrospective**

#### Final Presentation (45 minutes)
```
Agenda:

1. Product Demo (15 min):
   - Live demo of all features
   - Show real-time capabilities (location tracking)
   - Handle error scenarios gracefully

2. Technical Deep Dive (15 min):
   - Architecture walkthrough
   - Code highlights (interesting algorithms)
   - Performance metrics:
     * Load test results (k6 report)
     * Latency graphs
     * Database query performance
   - Observability:
     * Show Grafana dashboards
     * Demo log aggregation
     * Tracing example

3. Product Decisions (10 min):
   - Feature prioritization rationale
   - A/B test results (mock data OK)
   - Trade-offs made

4. Lessons Learned (5 min):
   - What went well?
   - What would you do differently?
   - What's next? (Future roadmap)

Evaluation Rubric:
┌─────────────────────────┬────────┬─────────────────────────┐
│ Criteria                │ Weight │ Evaluation              │
├─────────────────────────┼────────┼─────────────────────────┤
│ Code Quality            │ 25%    │ Clean, tested, reviewed │
│ System Design           │ 25%    │ Scalable, resilient     │
│ Product Thinking        │ 20%    │ User-focused, metrics   │
│ Observability           │ 15%    │ Logs, metrics, traces   │
│ Presentation            │ 15%    │ Clear, confident        │
└─────────────────────────┴────────┴─────────────────────────┘

Passing Score: 80%
```

#### Retrospective Session
```
Format: Start-Stop-Continue

What should we START doing?
- More pair programming sessions
- Weekly tech talks

What should we STOP doing?
- Long meetings without agenda
- Reviewing code without context

What should we CONTINUE doing?
- Daily standups
- Blameless post-mortems
- Code review culture

Action Items:
- [ ] Document common pitfalls for next cohort
- [ ] Update training materials based on feedback
- [ ] Create knowledge base of best practices
```

---

# 📊 EVALUATION & GRADUATION

## Performance Review Scorecard

### Technical Skills (60%)
```
Clean Code & Testing        ▓▓▓▓▓▓▓▓▓░ 90%
Git & Collaboration         ▓▓▓▓▓▓▓▓░░ 85%
Database & Optimization     ▓▓▓▓▓▓▓▓▓░ 88%
System Design               ▓▓▓▓▓▓▓░░░ 75%
Observability              ▓▓▓▓▓▓▓▓░░ 82%
CI/CD & Infra              ▓▓▓▓▓▓▓▓░░ 80%
```

### Collaboration Skills (25%)
```
Code Review Quality         ▓▓▓▓▓▓▓▓▓░ 92%
Communication              ▓▓▓▓▓▓▓▓░░ 85%
Ownership & Accountability  ▓▓▓▓▓▓▓▓▓░ 88%
```

### Product Mindset (15%)
```
User Empathy               ▓▓▓▓▓▓▓░░░ 78%
Data-Driven Decisions      ▓▓▓▓▓▓▓▓░░ 83%
```

## Graduation Criteria
✅ Overall Score: ≥ 80%
✅ No category < 70%
✅ Capstone project completed
✅ All weekly milestones achieved
✅ Positive peer reviews (≥ 4/5)

## Post-Graduation Path

### Level 1: Junior Engineer (Month 1-6)
- Work on well-defined tasks
- Pair with senior engineers
- Focus on execution speed

### Level 2: Mid-Level Engineer (Month 7-18)
- Own features end-to-end
- Mentor junior engineers
- Lead small projects

### Level 3: Senior Engineer (18+ months)
- Design complex systems
- Set technical direction
- Cross-team collaboration

---

# 🛠️ RESOURCES & TOOLS

## Required Reading
- [ ] Clean Code (Robert C. Martin)
- [ ] Designing Data-Intensive Applications (Martin Kleppmann)
- [ ] The Phoenix Project (Gene Kim)
- [ ] Accelerate (Nicole Forsgren)

## Tools Setup
```bash
# Development
- IDE: VSCode/IntelliJ/GoLand
- Git: GitHub/GitLab
- API Testing: Postman/Insomnia

# Infrastructure
- Docker Desktop
- Kubernetes: minikube/kind
- Cloud: AWS Free Tier/GCP

# Monitoring
- Prometheus + Grafana (local)
- ELK Stack (Elasticsearch, Logstash, Kibana)

# Load Testing
- k6
- Apache JMeter
```

## Slack Channels
- #training-general (announcements)
- #training-help (ask questions)
- #training-wins (celebrate achievements)
- #code-review (request reviews)

## Office Hours
- Mentor 1-on-1: Tuesday & Thursday 2-3pm
- Group Q&A: Friday 4-5pm
- Code Review Sessions: Daily 10-11am

---

# 🎓 ALUMNI TESTIMONIALS

> "Tuần đầu mình struggle với TDD, nhưng giờ không thể tưởng tượng code mà không có test. Game changer!" - Minh Nguyen, Software Engineer @ Grab

> "Capstone project giống như làm startup mini. Học được nhiều hơn cả năm đại học!" - Thu Tran, Backend Engineer @ Shopee

> "Code review culture ở đây brutal nhưng mình improve nhanh gấp 10 lần. Worth it!" - Huy Le, Full-Stack Engineer @ Tiki

---

**Good luck! Remember: Code is read 10x more than it's written. Write for the next person, not for the compiler. 🚀**
