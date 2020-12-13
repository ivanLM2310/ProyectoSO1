#include <linux/fs.h>
#include <linux/init.h>
#include <linux/kernel.h>
#include <linux/list.h>
#include <linux/module.h>
#include <linux/proc_fs.h>
#include <linux/sched.h>
#include <linux/seq_file.h>
#include <linux/slab.h>
#include <linux/string.h>
#include <linux/types.h>
#include <linux/mm.h>
#include <linux/sched/signal.h>


char buffer[256];
struct sysinfo info;
struct task_struct *task_list;
struct task_struct *task_child;
struct list_head *list;

char * get_task_state(long state)
{
    switch (state) {
        case TASK_RUNNING:
            return "Ejecucion";
        case TASK_INTERRUPTIBLE:
            return "Interrumpible";
        case TASK_UNINTERRUPTIBLE:
            return "Ininterrumpible";
        case __TASK_STOPPED:
            return "Detenido";
        case __TASK_TRACED:
            return "TASK_TRACED";
        case TASK_STOPPED:
            return "Detenido";
        case EXIT_ZOMBIE:
            return "Zombie";
        default:
        {
            sprintf(buffer, "Desconocido%ld\n", state);
            return buffer;
        }
    }
}

void write_process(struct seq_file *m, struct task_struct *s,long totalM){
    #define Convert(x) ((x) << (PAGE_SHIFT - 10))
    seq_printf(m,"PID: %d\t\tNombre: %s\t\tMemoria: %ld\t\tTotalM: %ld\t\tEstado: %s\n",s->pid, s->comm,get_mm_rss(s->mm), Convert(totalM), get_task_state(s->state));
    #undef K
}

void write_process_json(struct seq_file *m, struct task_struct *s,long totalM){
    #define Convert(x) ((x) << (PAGE_SHIFT - 10))
    seq_printf(m,"\"PID\":\"%d\",\"Nombre\":\"%s\",\"Memoria\":\"%ld\",\"TotalM\":\"%ld\",\"Estado\":\"%s\"\n",s->pid, s->comm,get_mm_rss(s->mm), Convert(totalM), get_task_state(s->state));
    #undef K
}

static int pstreeG(struct seq_file *m, void *v)
{
    
    unsigned int process_count = 0;
    int i = 0;
    int j = 0;

    pr_info("%s: In init\n", __func__);
    si_meminfo(&info);  
    seq_printf(m,"[\n");
    for_each_process(task_list) {
        if(task_list->mm){
            if(i==0){
                i +=1;
            }else{
                seq_printf(m,",");
            }
            seq_printf(m,"{\n");
            write_process_json(m,task_list,info.totalram);
            
            j = 0;
            list_for_each(list, &task_list->children){   
                if(j == 0){
                    seq_printf(m,",\"hijos\":[\n");
                    j+=1;
                }else{
                    seq_printf(m,",");
                }
                seq_printf(m,"{\n");                
                task_child = list_entry( list, struct task_struct, sibling );   
                write_process_json(m,task_child,info.totalram);
                seq_printf(m,"}\n");

            }
            if(j > 0){
                seq_printf(m,"]\n");
            }
            
            seq_printf(m,"}\n");
            
        }
        process_count++;    
    }
    seq_printf(m,"]\n");
    pr_info("Number of processes:%u\n", process_count);

    return 0;
}

static int cpu_info_proc_open(struct inode *inode, struct file *file)
{
    return single_open(file, pstreeG, NULL);
}

static const struct file_operations cpu_info_proc_fops = {
    .open       = cpu_info_proc_open,
    .read       = seq_read,
    .llseek     = seq_lseek,
    .release    = single_release,
};

MODULE_LICENSE("GPL");
MODULE_DESCRIPTION("Modulo de CPU - Sistemas Operativos 2");

static int __init cpu_grupo22_init(void)
{
	printk(KERN_INFO "Hola mundo, somos el grupo 22(Los Cracks :v) y este es el monitor de CPU\n");
	proc_create("cpu_grupo22", 0, NULL, &cpu_info_proc_fops);
	return 0;
}

static void __exit cpu_grupo22_cleanup(void)
{
	remove_proc_entry("cpu_grupo22", NULL);
	printk(KERN_INFO "Sayonara mundo, somos el grupo 22(Los Cracks :v) y este fue el monitor de CPU\n");
}

module_init(cpu_grupo22_init);
module_exit(cpu_grupo22_cleanup);