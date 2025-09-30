# Product Microservice - Documentation Index

Welcome! This document helps you navigate all the documentation for this project.

## 🚀 Getting Started

**New to this project? Start here:**

1. **[QUICKSTART.md](QUICKSTART.md)** - Get up and running in 5 minutes
   - Prerequisites check
   - Installation steps
   - First API call
   - Common issues and solutions

2. **[README.md](README.md)** - Comprehensive documentation
   - Full feature list
   - Detailed API documentation
   - Testing guide
   - Documentation references used

## 📚 Understanding the Project

**Want to understand how it works?**

3. **[PROJECT_SUMMARY.md](PROJECT_SUMMARY.md)** - High-level overview
   - Technical stack
   - Core functionality
   - Key implementation details
   - Lessons learned

4. **[ARCHITECTURE.md](ARCHITECTURE.md)** - Deep dive into architecture
   - System architecture diagrams
   - Layer responsibilities
   - Data flow explanations
   - Database relationships
   - Error handling strategy

## ✅ Verification and Submission

**Ready to verify everything works?**

5. **[SUBMISSION_CHECKLIST.md](SUBMISSION_CHECKLIST.md)** - Complete requirements checklist
   - All requirements verified
   - Documentation references
   - Evaluation criteria met
   - Quick verification steps

## 🐛 Learning Materials

**Want to learn from common mistakes?**

6. **[debug_exercise/DEBUG_EXPLANATION.md](debug_exercise/DEBUG_EXPLANATION.md)** - Debug exercise
   - Common UUID and GORM issues
   - Documentation used to solve problems
   - Step-by-step solutions
   - Key takeaways

## 📁 Code Organization

### Main Application
```
cmd/server/main.go          - Application entry point
```

### Core Implementation
```
internal/
├── database/               - Database setup and migrations
├── handler/               - gRPC handlers (presentation layer)
│   ├── product_handler.go
│   ├── product_handler_test.go
│   ├── subscription_handler.go
│   └── subscription_handler_test.go
├── models/                - Database models
│   ├── product.go
│   └── subscription_plan.go
├── repository/            - Data access layer
│   ├── product_repository.go
│   ├── product_repository_test.go
│   ├── subscription_repository.go
│   └── subscription_repository_test.go
└── service/               - Business logic layer
    ├── product_service.go
    ├── product_service_test.go
    ├── subscription_service.go
    └── subscription_service_test.go
```

### API Definitions
```
proto/
├── product.proto          - Product service definition
└── subscription.proto     - Subscription service definition
```

### Examples and Scripts
```
examples/client/main.go    - Example Go client
scripts/
├── setup.sh              - Setup script (Linux/macOS)
├── test_api.sh           - API test script (Linux/macOS)
└── test_api.bat          - API test script (Windows)
```

## 🎯 Quick Reference

### Common Commands

| Task | Command |
|------|---------|
| Setup environment | `bash scripts/setup.sh` |
| Generate proto files | `make proto` |
| Run tests | `make test` |
| Start server | `make run` |
| Test API | `bash scripts/test_api.sh` |
| Run example client | `go run examples/client/main.go` |
| Generate coverage | `make test-coverage` |

### API Endpoints

**ProductService** (5 endpoints):
- CreateProduct
- GetProduct
- UpdateProduct
- DeleteProduct
- ListProducts

**SubscriptionService** (5 endpoints):
- CreateSubscriptionPlan
- GetSubscriptionPlan
- UpdateSubscriptionPlan
- DeleteSubscriptionPlan
- ListSubscriptionPlans

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| DB_DRIVER | sqlite | Database driver |
| DB_HOST | localhost | Database host |
| DB_PORT | 5432 | Database port |
| DB_NAME | products.db | Database name |
| PORT | 50051 | gRPC server port |

## 📖 Documentation by Use Case

### "I want to understand the requirements"
→ Read [SUBMISSION_CHECKLIST.md](SUBMISSION_CHECKLIST.md)

### "I want to run the project quickly"
→ Follow [QUICKSTART.md](QUICKSTART.md)

### "I want to understand the architecture"
→ Read [ARCHITECTURE.md](ARCHITECTURE.md)

### "I want to see API examples"
→ Check [README.md](README.md) API Documentation section

### "I want to learn about common issues"
→ Study [debug_exercise/DEBUG_EXPLANATION.md](debug_exercise/DEBUG_EXPLANATION.md)

### "I want to see working code examples"
→ Run `examples/client/main.go`

### "I want to understand the design decisions"
→ Read [PROJECT_SUMMARY.md](PROJECT_SUMMARY.md)

### "I want to test the API"
→ Use `scripts/test_api.sh` or `scripts/test_api.bat`

## 🔍 Finding Specific Information

### Database Related
- **Schema**: [ARCHITECTURE.md](ARCHITECTURE.md) - Database Relationships section
- **Migrations**: [README.md](README.md) - Data Storage section
- **Models**: `internal/models/` directory

### gRPC Related
- **Proto files**: `proto/` directory
- **Handlers**: `internal/handler/` directory
- **Examples**: [README.md](README.md) - API Documentation section

### Testing Related
- **Test files**: `*_test.go` files in each package
- **Running tests**: [README.md](README.md) - Testing section
- **Test strategy**: [ARCHITECTURE.md](ARCHITECTURE.md) - Testing Architecture section

### Configuration Related
- **Environment variables**: [README.md](README.md) - Environment Variables section
- **Database config**: `internal/database/database.go`
- **Example config**: `.env.example`

## 📚 External Documentation References

All external documentation used is listed in:
- [README.md](README.md) - "Documentation Used" section
- [debug_exercise/DEBUG_EXPLANATION.md](debug_exercise/DEBUG_EXPLANATION.md) - "Documentation Used" section

Key resources:
1. gRPC Go Documentation
2. GORM Documentation
3. Protocol Buffers Guide
4. Go Project Layout Standards
5. Testify Documentation

## 🎓 Learning Path

**Recommended order for learning:**

1. Start with [QUICKSTART.md](QUICKSTART.md) to get it running
2. Read [PROJECT_SUMMARY.md](PROJECT_SUMMARY.md) for overview
3. Study [ARCHITECTURE.md](ARCHITECTURE.md) for deep understanding
4. Review code in this order:
   - `internal/models/` - Understand data structures
   - `internal/repository/` - See data access patterns
   - `internal/service/` - Learn business logic
   - `internal/handler/` - Understand gRPC integration
5. Read [debug_exercise/DEBUG_EXPLANATION.md](debug_exercise/DEBUG_EXPLANATION.md) to learn from mistakes
6. Run tests to see everything in action
7. Experiment with the example client

## 🆘 Getting Help

### Troubleshooting
1. Check [QUICKSTART.md](QUICKSTART.md) - Common Issues section
2. Review [README.md](README.md) - Common Issues and Solutions section
3. Look at test files for usage examples

### Understanding Errors
1. Check [debug_exercise/DEBUG_EXPLANATION.md](debug_exercise/DEBUG_EXPLANATION.md)
2. Review error handling in [ARCHITECTURE.md](ARCHITECTURE.md)

### API Usage
1. See examples in [README.md](README.md) - API Documentation section
2. Run `examples/client/main.go` for working code
3. Use `scripts/test_api.sh` for command-line examples

## 📊 Project Statistics

- **Total Files**: 30+ source files
- **Lines of Code**: ~3000+ lines
- **Test Coverage**: Service, Repository, and Handler layers
- **Documentation**: 7 comprehensive markdown files
- **API Endpoints**: 10 gRPC methods
- **Database Tables**: 2 with relationships

## ✨ Highlights

- ✅ Complete implementation of all requirements
- ✅ Clean architecture with clear separation of concerns
- ✅ Comprehensive testing at all layers
- ✅ Extensive documentation with examples
- ✅ Production-ready code quality
- ✅ Debug exercise demonstrating learning
- ✅ Multiple ways to test and interact with the API

## 🎯 Next Steps

After reviewing the documentation:

1. **Run the project**: Follow [QUICKSTART.md](QUICKSTART.md)
2. **Explore the code**: Start with `cmd/server/main.go`
3. **Run tests**: `make test`
4. **Try the API**: Use grpcurl or the example client
5. **Read the debug exercise**: Learn from common mistakes
6. **Extend the project**: Add new features using the same patterns

---

**Happy coding! 🚀**

For any questions, refer to the specific documentation file listed above.

