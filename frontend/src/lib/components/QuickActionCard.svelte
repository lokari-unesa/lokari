<script lang="ts">
  import { ChevronRight } from "lucide-svelte";
  import { cn } from "$lib/utils";

  let { to, icon: Icon, title, description, accent = "primary" } = $props<{
    to: string,
    icon: any,
    title: string,
    description: string,
    accent?: "primary" | "secondary" | "warning" | "safe"
  }>();

  const accents: Record<string, string> = {
    primary: "bg-primary/10 text-primary",
    secondary: "bg-secondary/10 text-secondary",
    warning: "bg-warning/10 text-warning",
    safe: "bg-safe/10 text-safe",
  };
</script>

<a
  href={to}
  class="group relative flex flex-col p-5 rounded-xl bg-card border border-border hover:border-primary/40 hover:shadow-[0_4px_12px_-2px_rgba(15,23,42,0.08)] transition-all"
>
  <span class={cn("grid place-items-center w-11 h-11 rounded-xl", accents[accent])}>
    <Icon class="w-6 h-6" strokeWidth={2} />
  </span>
  <h3 class="mt-4 font-display font-semibold text-[1.125rem] text-foreground leading-tight">{title}</h3>
  <p class="mt-1.5 text-[0.9375rem] text-muted-foreground leading-relaxed">{description}</p>
  <span class="mt-4 inline-flex items-center gap-1 text-[0.8125rem] font-semibold text-primary">
    Buka <ChevronRight class="w-3.5 h-3.5 group-hover:translate-x-0.5 transition-transform" />
  </span>
</a>
