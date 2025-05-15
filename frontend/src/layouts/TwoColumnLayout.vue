<template>
  <div class="min-h-screen bg-nordic-primary-bg">
    <div class="container mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div class="flex flex-col md:flex-row md:space-x-8">
        <aside class="w-full md:w-64 mb-6 md:mb-0 bg-white p-6 rounded-lg shadow-sm">
          <slot name="sidebar"></slot>
        </aside>

        <main class="w-full md:flex-1">
          <slot name="main"></slot>
        </main>
      </div>
    </div>
  </div>
</template>

<script setup>
// No props or reactive state needed for this simple layout
// The setup script tag automatically exposes everything to the template
</script>

<style scoped>
.two-column-layout {
  display: flex;
  width: 100%;
  min-height: 100vh;
}

.sidebar {
  width: 280px;
  padding: 20px;
  background-color: #f5f5f5;
  border-right: 1px solid #e0e0e0;
}

.main-content {
  flex: 1;
  padding: 20px;
  background-color: #ffffff;
}

/* Responsive behavior */
@media (max-width: 768px) {
  .two-column-layout {
    flex-direction: column;
  }

  .sidebar {
    width: 100%;
    border-right: none;
    border-bottom: 1px solid #e0e0e0;
  }
}
</style>